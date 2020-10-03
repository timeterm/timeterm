package authn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	microsoftoauth2 "golang.org/x/oauth2/microsoft"

	"gitlab.com/timeterm/timeterm/backend/database"
)

type Issuer struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

type Authorizer struct {
	dbw     *database.Wrapper
	log     logr.Logger
	issuers map[string]Issuer

	redirectURL *url.URL
}

func New(dbw *database.Wrapper, log logr.Logger) (*Authorizer, error) {
	redirectURL, err := url.Parse(os.Getenv("OIDC_REDIRECT_URL"))
	if err != nil {
		return nil, fmt.Errorf("invalid OIDC_REDIRECT_URL: %w", err)
	}

	a := &Authorizer{
		dbw:         dbw,
		log:         log,
		issuers:     make(map[string]Issuer),
		redirectURL: redirectURL,
	}

	err = a.setupIssuers()
	if err != nil {
		return nil, fmt.Errorf("could not set up issuers: %w", err)
	}
	return a, nil
}

func (a *Authorizer) RegisterRoutes(r *echo.Echo) {
	r.GET("/oidc/login/:issuer", a.HandleLogin)
	r.GET("/oidc/callback", a.HandleOauth2Callback)
}

func (a *Authorizer) setupIssuers() error {
	err := a.setupIssuerGoogle()
	if err != nil {
		return fmt.Errorf("could not setup Google issuer: %w", err)
	}

	err = a.setupIssuerMicrosoft()
	if err != nil {
		return fmt.Errorf("could not setup Microsoft issuer: %w", err)
	}

	return nil
}

func (a *Authorizer) setupIssuerGoogle() error {
	googleClientID := os.Getenv("OIDC_PROVIDERS_GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("OIDC_PROVIDERS_GOOGLE_CLIENT_SECRET")
	googleProvider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("could not setup OIDC provider Google: %w", err)
	}

	a.issuers["google"] = Issuer{
		config: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  a.redirectURL.String(),
			Endpoint:     googleProvider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: googleProvider.Verifier(&oidc.Config{
			ClientID: googleClientID,
		}),
	}

	return nil
}

const microsoftIssuerURL = "https://login.microsoftonline.com/common/v2.0"
const microsoftJWKSURL = "https://login.microsoftonline.com/common/discovery/v2.0/keys"

func (a *Authorizer) setupIssuerMicrosoft() error {
	microsoftClientID := os.Getenv("OIDC_PROVIDERS_MICROSOFT_CLIENT_ID")
	microsoftClientSecret := os.Getenv("OIDC_PROVIDERS_MICROSOFT_CLIENT_SECRET")

	keySet := oidc.NewRemoteKeySet(context.Background(), microsoftJWKSURL)

	a.issuers["microsoft"] = Issuer{
		config: &oauth2.Config{
			ClientID:     microsoftClientID,
			ClientSecret: microsoftClientSecret,
			RedirectURL:  a.redirectURL.String(),
			Endpoint:     microsoftoauth2.AzureADEndpoint("common"),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: oidc.NewVerifier(microsoftIssuerURL, keySet, &oidc.Config{
			ClientID:             microsoftClientID,
			SkipIssuerCheck:      true,
			SupportedSigningAlgs: []string{oidc.RS256},
		}),
	}

	return nil
}

func (a *Authorizer) HandleLogin(c echo.Context) error {
	redirectURL := c.QueryParam("redirectTo")

	issuerName := c.Param("issuer")
	issuer, ok := a.issuers[issuerName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown issuer")
	}

	state, err := a.dbw.CreateOAuth2State(c.Request().Context(), issuerName, redirectURL)
	if err != nil {
		a.log.Error(err, "could not create token")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create token")
	}

	return c.Redirect(http.StatusFound, issuer.config.AuthCodeURL(state.State.String()))
}

func (a *Authorizer) issuerFromState(state database.OAuth2State) (Issuer, error) {
	issuer, ok := a.issuers[state.Issuer]
	if !ok {
		return Issuer{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid issuer")
	}
	return issuer, nil
}

func (a *Authorizer) stateFromRequest(ctx context.Context, data oauth2CallbackRequest) (database.OAuth2State, error) {
	state, err := uuid.Parse(data.State)
	if err != nil {
		return database.OAuth2State{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid state")
	}

	stateInfo, err := a.dbw.GetOAuth2State(ctx, state)
	if err != nil {
		return database.OAuth2State{}, echo.NewHTTPError(http.StatusBadRequest, "Nonexistent state")
	}

	return stateInfo, nil
}

type oauth2CallbackRequest struct {
	State string `form:"state"`
	Code  string `query:"code"`
}

func (a *Authorizer) HandleOauth2Callback(c echo.Context) error {
	ctx := context.Background()

	var reqData oauth2CallbackRequest
	if err := c.Bind(&reqData); err != nil {
		return err
	}

	state, err := a.stateFromRequest(ctx, reqData)
	if err != nil {
		return err
	}

	issuer, err := a.issuerFromState(state)
	if err != nil {
		return err
	}

	oauth2Token, err := issuer.config.Exchange(ctx, reqData.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Could not exchange code with provider",
		)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusPreconditionFailed,
			"Could not exchange code with provider (id_token not present or string)",
		)
	}

	idToken, err := issuer.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not verify token")
	}

	var claims struct {
		Issuer     string `json:"iss"`
		Subject    string `json:"sub"`
		Audience   string `json:"aud"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Verified   bool   `json:"email_verified"`
		Locale     string `json:"locale"` // BCP 47
		Picture    string `json:"picture"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not read claims")
	}

	user, err := a.dbw.GetUserByEmail(c.Request().Context(), claims.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			a.log.Error(err, "could not get user by email")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
		}
	}
	userExists := err == nil

	user, err = a.dbw.GetUserByOIDCFederation(c.Request().Context(), database.OIDCFederation{
		OIDCIssuer:  claims.Issuer,
		OIDCSubject: claims.Subject,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			a.log.Error(err, "could not get user by OIDC federation")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
		}
	}
	userPreviouslyLoggedInWithIssuer := err == nil

	if userExists {
		if !userPreviouslyLoggedInWithIssuer {
			// TODO(rutgerbrf): do this some other way for the user to be able to use multiple providers.
			return echo.NewHTTPError(http.StatusBadRequest, "User created with different provider")
		}
	} else if !userPreviouslyLoggedInWithIssuer {
		org, err := a.dbw.CreateOrganization(context.Background(), "", "")
		if err != nil {
			a.log.Error(err, "could not create organization")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create organization")
		}

		// TODO(rutgerbrf): do CreateUser & CreateOIDCFederation in a transaction to prevent corruption.
		user, err = a.dbw.CreateUser(context.Background(), claims.Name, claims.Email, org.ID)
		if err != nil {
			a.log.Error(err, "could not create user")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create user")
		}

		_, err = a.dbw.CreateOIDCFederation(context.Background(), database.OIDCFederation{
			OIDCIssuer:   claims.Issuer,
			OIDCSubject:  claims.Subject,
			OIDCAudience: claims.Audience,
			UserID:       user.ID,
		})
		if err != nil {
			a.log.Error(err, "could not create OIDC federation")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not save login info")
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "No user for login")
	}

	token, err := a.dbw.CreateToken(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create token")
	}

	c.SetCookie(&http.Cookie{
		Name:    "ttsess",
		Value:   token.String(),
		Expires: time.Now().Add(time.Hour * 24),
	})

	return c.Redirect(http.StatusSeeOther, state.RedirectURL)
}

package authn

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	microsoftoauth2 "golang.org/x/oauth2/microsoft"

	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/templates"
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

func (a *Authorizer) RegisterRoutes(g *echo.Group) {
	g = g.Group("/oidc")
	g.GET("/login/:issuer", a.HandleLogin)
	g.GET("/callback", a.HandleOauth2Callback)
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

type Status string

const (
	StatusError Status = "error"
	StatusOK    Status = "ok"
)

func errorMsg(msg string) url.Values {
	return url.Values{"error": []string{msg}}
}

func tokenData(token string) url.Values {
	return url.Values{
		"token": []string{
			base64.URLEncoding.EncodeToString([]byte(token)),
		},
		"expires": []string{
			strconv.FormatInt(time.Now().Add(database.DefaultTokenExpiration).Unix(), 10),
		},
	}
}

func redirectToOrigin(c echo.Context, redirectTo *url.URL, status Status, data url.Values) error {
	origin := *redirectTo

	q := origin.Query()
	q.Set("status", string(status))
	if data != nil {
		for k, v := range data {
			q[k] = v
		}
	}
	origin.RawQuery = q.Encode()

	return c.Redirect(http.StatusFound, origin.String())
}

func (a *Authorizer) HandleLogin(c echo.Context) error {
	redirectURL, err := url.Parse(c.QueryParam("redirectTo"))
	if err != nil {
		return err
	}

	issuerName := c.Param("issuer")
	issuer, ok := a.issuers[issuerName]
	if !ok {
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Unknown issuer"))
	}

	state, err := a.dbw.CreateOAuth2State(c.Request().Context(), issuerName, redirectURL.String())
	if err != nil {
		a.log.Error(err, "could not create token")
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not create token"))
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

func (a *Authorizer) stateFromRequest(c echo.Context, data oauth2CallbackRequest) (database.OAuth2State, error) {
	state, err := uuid.Parse(data.State)
	if err != nil {
		return database.OAuth2State{}, templates.RenderError(c, http.StatusBadRequest, "Invalid state")
	}

	stateInfo, err := a.dbw.GetOAuth2State(c.Request().Context(), state)
	if err != nil {
		return database.OAuth2State{}, templates.RenderError(c, http.StatusBadRequest, "Nonexistent state")
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

	state, err := a.stateFromRequest(c, reqData)
	if err != nil {
		return err
	}

	redirectURL, err := url.Parse(state.RedirectURL)
	if err != nil {
		return templates.RenderError(c, http.StatusInternalServerError, "Invalid redirect URL")
	}

	issuer, err := a.issuerFromState(state)
	if err != nil {
		return err
	}

	oauth2Token, err := issuer.config.Exchange(ctx, reqData.Code)
	if err != nil {
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not exchange code with provider"))
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not exchange code with provider"))
	}

	idToken, err := issuer.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not verify token"))
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
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not read claims"))
	}

	user, err := a.dbw.GetUserByEmail(c.Request().Context(), claims.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			a.log.Error(err, "could not get user by email")
			return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not query database"))
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
			return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not query database"))
		}
	}
	userPreviouslyLoggedInWithIssuer := err == nil

	if userExists {
		if !userPreviouslyLoggedInWithIssuer {
			// TODO(rutgerbrf): do this some other way for the user to be able to use multiple providers.
			return redirectToOrigin(c, redirectURL, StatusError, errorMsg("User created with different provider"))
		}
	} else if !userPreviouslyLoggedInWithIssuer {
		user, err = a.dbw.CreateNewUser(context.Background(), claims.Name, claims.Email, database.OIDCFederation{
			OIDCIssuer:   claims.Issuer,
			OIDCSubject:  claims.Subject,
			OIDCAudience: claims.Audience,
		})
		if err != nil {
			a.log.Error(err, "could not create user")
			return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not create user"))
		}
	} else {
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("No user for login"))
	}

	token, err := a.dbw.CreateUserToken(c.Request().Context(), user.ID)
	if err != nil {
		a.log.Error(err, "could not create token")
		return redirectToOrigin(c, redirectURL, StatusError, errorMsg("Could not create token"))
	}

	return redirectToOrigin(c, redirectURL, StatusOK, tokenData(token.String()))
}

package authn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

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
	issuerName := c.Param("issuer")
	issuer, ok := a.issuers[issuerName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown issuer")
	}

	state, err := a.dbw.CreateOAuth2State(c.Request().Context(), issuerName, a.redirectURL.String())
	if err != nil {
		a.log.Error(err, "could not create token")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create token")
	}

	return c.Redirect(http.StatusFound, issuer.config.AuthCodeURL(state.State.String()))
}

func (a *Authorizer) issuerNameFromState(ctx context.Context, state uuid.UUID) (string, error) {
	stateInfo, err := a.dbw.GetOAuth2State(ctx, state)
	if err != nil {
		return "", fmt.Errorf("could not get OAuth2 state: %w", err)
	}

	return stateInfo.Issuer, nil
}

func (a *Authorizer) issuerFromState(ctx context.Context, state uuid.UUID) (Issuer, error) {
	issuerName, err := a.issuerNameFromState(ctx, state)
	if err != nil {
		return Issuer{}, err
	}

	issuer, ok := a.issuers[issuerName]
	if !ok {
		return Issuer{}, errors.New("unknown issuer")
	}
	return issuer, nil
}

func (a *Authorizer) issuerFromRequest(ctx context.Context, data oauth2CallbackData) (Issuer, error) {
	state, err := uuid.Parse(data.State)
	if err != nil {
		return Issuer{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid state")
	}

	issuer, err := a.issuerFromState(ctx, state)
	if err != nil {
		return Issuer{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid issuer")
	}
	return issuer, nil
}

type oauth2CallbackData struct {
	State string `form:"state"`
	Code  string `query:"code"`
}

func (a *Authorizer) HandleOauth2Callback(c echo.Context) error {
	var data oauth2CallbackData
	if err := c.Bind(&data); err != nil {
		return err
	}

	issuer, err := a.issuerFromRequest(c.Request().Context(), data)
	if err != nil {
		return err
	}

	oauth2Token, err := issuer.config.Exchange(c.Request().Context(), data.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not exchange code with provider")
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Could not exchange code with provider (id_token not present or string)")
	}

	idToken, err := issuer.verifier.Verify(c.Request().Context(), rawIDToken)
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

	a.log.Info("user authenticated",
		"email", claims.Email,
		"subject", claims.Subject,
		"issuer", claims.Issuer,
		"audience", claims.Audience,
	)

	return nil
}

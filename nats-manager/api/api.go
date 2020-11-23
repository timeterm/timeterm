// Package api implements a very small subset of the nats-account-server API.
// Using nats-account-server would simply be cumbersome for our use case,
// and the API that is provides is relatively trivial to implement.
package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/nats-io/jwt/v2"

	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/vla"
	vlahttprouter "gitlab.com/timeterm/timeterm/nats-manager/pkg/vla/impl/httprouter"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

type Server struct {
	r       vla.Router
	log     logr.Logger
	secrets *secrets.Store
	mgr     *manager.Manager
}

func NewServer(log logr.Logger, store *secrets.Store, mgr *manager.Manager) *Server {
	s := Server{
		r:       vlahttprouter.New(),
		log:     log.WithName("ApiServer"),
		secrets: store,
		mgr:     mgr,
	}
	s.registerRoutes()

	return &s
}

func (s *Server) Serve(ctx context.Context, addr string) error {
	const shutdownTimeout = time.Second * 30

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 30,
		Handler:      s.r,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	return srv.ListenAndServe()
}

func (s *Server) registerRoutes() {
	vla.GET(s.r, "/jwt/v1/accounts/", func(w http.ResponseWriter, _ *http.Request, _ vla.Route, _ vla.Params) {
		h := w.Header()
		h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Set("Pragma", "no-cache")
		h.Set("Expires", "0")
		w.WriteHeader(http.StatusOK)
	})
	vla.GET(s.r, "/jwt/v1/accounts/:pubkey", s.GetJWT)
	vla.GET(s.r, "/creds/v1/accounts/:account/users/:user/", s.GetUserCreds)
}

func (s *Server) GetJWT(w http.ResponseWriter, r *http.Request, _ vla.Route, p vla.Params) {
	token, err := s.secrets.ReadJWTLiteral(p.ByName("pubkey"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Must be an account
	claims, err := jwt.DecodeAccountClaims(token)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.URL.Query().Get("check") == "true" {
		var vr jwt.ValidationResults
		claims.Validate(&vr)

		if len(vr.Errors()) > 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	etag := strconv.Quote(claims.Claims().ID)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	w.Header().Set("ETag", etag)

	w.Header().Set("Content-Type", "application/jwt")
	if r.URL.Query().Get("text") == "true" {
		w.Header().Set("Content-Type", "text/plain")
	}

	if expiresIn := claims.Claims().Expires - time.Now().Unix(); expiresIn > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", expiresIn))
	}

	if r.URL.Query().Get("decode") == "true" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(claims.String()))
		claims.Payload()
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(token))
}

func (s *Server) GetUserCreds(w http.ResponseWriter, r *http.Request, _ vla.Route, p vla.Params) {
	account := p.ByName("account")
	user := p.ByName("user")

	creds, err := s.mgr.GenerateUserCredentials(r.Context(), user, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write([]byte(creds))
}

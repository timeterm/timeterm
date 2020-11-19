// Package api implements a very small subset of the nats-account-server API.
// Using nats-account-server would simply be cumbersome for our use case,
// and the API that is provides is relatively trivial to implement.
package api

import (
	"net/http"

	"gitlab.com/timeterm/timeterm/nats-manager/pkg/vla"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

type Server struct {
	r     vla.Router
	vault *secrets.VaultClient
}

func (s *Server) registerRoutes() {
	s.r.GET("/jwt/v1/accounts/:id", s.GetJWT)
}

func (s *Server) GetJWT(w http.ResponseWriter, r *http.Request, vr vla.Route, p vla.Params) {
	s.vault.ReadU
}

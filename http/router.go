package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func Router(h Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(h.tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Post("/users/signup", h.SignUp)

	return r
}

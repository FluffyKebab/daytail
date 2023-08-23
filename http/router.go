package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func router(h Handler) http.Handler {
	r := chi.NewRouter()

	// Protected routes.
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/users/{userId}/entries", h.createEntry)
		r.Get("/users/{userId}/entries", h.getEntries)
	})

	// Public routes.
	r.Group(func(r chi.Router) {
		r.Post("/users", h.signUp)
	})

	return r
}

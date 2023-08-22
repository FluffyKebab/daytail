package http

import (
	"net/http"

	"github.com/FluffyKebab/daytail"
	"github.com/go-chi/jwtauth/v5"
)

type Handler struct {
	daytail.EntryService
	daytail.UserService
	tokenAuth *jwtauth.JWTAuth
}

func ListenAndServe(address string, h Handler) error {
	return http.ListenAndServe(address, router(h))
}

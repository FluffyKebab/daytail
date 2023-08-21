package http

import (
	"github.com/FluffyKebab/daytail"
	"github.com/go-chi/jwtauth/v5"
)

type Handler struct {
	daytail.EntryService
	daytail.UserService

	tokenAuth *jwtauth.JWTAuth
}

func (h *Handler) InitAuth(jwtSecretKey string) {
	h.tokenAuth = jwtauth.New("HS256", []byte(jwtSecretKey), nil)
}

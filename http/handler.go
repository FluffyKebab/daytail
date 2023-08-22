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

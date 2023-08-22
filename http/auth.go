package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (h *Handler) InitAuth(jwtSecretKey string) {
	h.tokenAuth = jwtauth.New("HS256", []byte(jwtSecretKey), nil)
}

func (h *Handler) createToken(userId int) (string, error) {
	_, token, err := h.tokenAuth.Encode(map[string]interface{}{"userId": userId})
	return token, err
}

func (h *Handler) getUserId(r *http.Request) (int, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return 0, err
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("invalid user id in token")
	}

	return int(userId), nil
}

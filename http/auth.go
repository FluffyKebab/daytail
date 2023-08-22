package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func (h *Handler) InitAuth(jwtSecretKey string) {
	h.tokenAuth = jwtauth.New("HS256", []byte(jwtSecretKey), nil)
}

// authenticate tries to get a user token form the "userId" url param, and the
// jwt token and returns and the token if they are equal.
func (h *Handler) authenticate(r *http.Request) (int, error) {
	urlId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		return 0, errors.New("invalid or missing user id in url")
	}

	tokenId, err := h.getUserId(r)
	if err != nil {
		return 0, err
	}

	if urlId != tokenId {
		return 0, errors.New("unauthorized")
	}

	return urlId, nil
}

func (h *Handler) createToken(userId int) (string, error) {
	_, token, err := h.tokenAuth.Encode(map[string]interface{}{"userId": userId})
	return token, err
}

// getUserId gets a jwt token from the request.
func (h *Handler) getUserId(r *http.Request) (int, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return 0, err
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("invalid or missing user id in token")
	}

	return int(userId), nil
}

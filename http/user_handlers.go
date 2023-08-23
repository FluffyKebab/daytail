package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FluffyKebab/daytail"
)

type signUpRequestPayload struct {
	Name string `json:"name"`
}

type signUpResponsePayload struct {
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var u signUpRequestPayload
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = validateSignUpRequest(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.UserService.CreateUser(daytail.User{Name: u.Name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.createToken(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signUpResponsePayload{
		Token:  token,
		UserId: userId,
	})
}

func validateSignUpRequest(r signUpRequestPayload) error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

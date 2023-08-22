package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/FluffyKebab/daytail"
	"github.com/go-chi/chi/v5"
)

type createEntryRequestPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type createEntryResponsePayload struct {
	Id int `json:"id"`
}

// CreateEntry is the http handler for POST "/users/{userId}/entries".
func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var rp createEntryRequestPayload
	err := json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlParamUserId, err := validateCreateEntryRequest(rp, chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenUserId, err := h.getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if urlParamUserId != tokenUserId {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	entryId, err := h.EntryService.CreateEntry(daytail.Entry{
		UserID: urlParamUserId,
		Title:  rp.Title,
		Text:   rp.Text,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createEntryResponsePayload{Id: entryId})
}

func validateCreateEntryRequest(rp createEntryRequestPayload, userId string) (int, error) {
	if rp.Text == "" {
		return 0, errors.New("missing text")
	}

	if rp.Title == "" {
		return 0, errors.New("missing title")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return 0, errors.New("invalid or missing user id")
	}

	return id, nil
}

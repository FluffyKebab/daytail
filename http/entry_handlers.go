package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FluffyKebab/daytail"
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
	userId, err := h.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	var rp createEntryRequestPayload
	err = json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateCreateEntryRequest(rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entryId, err := h.EntryService.CreateEntry(daytail.Entry{
		UserID: userId,
		Title:  rp.Title,
		Text:   rp.Text,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createEntryResponsePayload{Id: entryId})
}

func validateCreateEntryRequest(rp createEntryRequestPayload) error {
	if rp.Text == "" {
		return errors.New("missing text")
	}

	if rp.Title == "" {
		return errors.New("missing title")
	}

	return nil
}

// GetEntries is the http handler for GET "/users/{userId}/entries".
func (h *Handler) GetEntries(w http.ResponseWriter, r *http.Request) {
	userId, err := h.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	entries, err := h.EntryService.UserEntries(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

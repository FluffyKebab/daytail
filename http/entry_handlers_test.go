package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FluffyKebab/daytail"
	"github.com/FluffyKebab/daytail/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	workingPayload, err := json.Marshal(createEntryRequestPayload{
		Title: "day 1",
		Text:  "foo baz",
	})
	require.NoError(t, err)

	var h Handler
	h.InitAuth("secret")
	workingToken, err := h.createToken(3)
	require.NoError(t, err)

	testCases := []struct {
		name                         string
		token                        string
		url                          string
		body                         io.Reader
		mock                         *mock.EntryService
		responseValidator            func(*httptest.ResponseRecorder)
		expectedCreateEntryIsInvoked bool
		expectedStatusCode           int
	}{
		{
			name:  "with valid payload",
			token: workingToken,
			url:   "/users/3/entries",
			body:  bytes.NewBuffer(workingPayload),
			mock: &mock.EntryService{
				CreateEntryFunc: func(d daytail.Entry) (int, error) {
					require.Equal(t, "day 1", d.Title)
					require.Equal(t, "foo baz", d.Text)
					return 1, nil
				},
			},
			responseValidator: func(rr *httptest.ResponseRecorder) {
				var responsePayload createEntryResponsePayload
				err := json.Unmarshal(rr.Body.Bytes(), &responsePayload)
				require.NoError(t, err)

				require.Equal(t, 1, responsePayload.Id)

			},
			expectedCreateEntryIsInvoked: true,
			expectedStatusCode:           http.StatusOK,
		},
		{
			name:  "unauthorized",
			token: workingToken,
			url:   "/users/5/entries",
			body:  bytes.NewBuffer(workingPayload),
			mock: &mock.EntryService{CreateEntryFunc: func(d daytail.Entry) (int, error) {
				return 0, nil
			}},
			responseValidator:            func(rr *httptest.ResponseRecorder) {},
			expectedCreateEntryIsInvoked: false,
			expectedStatusCode:           http.StatusUnauthorized,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", tt.url, tt.body)
			req.Header.Set("Authorization", "BEARER "+tt.token)

			h.EntryService = tt.mock
			Router(h).ServeHTTP(rw, req)

			require.Equal(t, tt.expectedStatusCode, rw.Result().StatusCode,
				"unexpected status code")
			require.Equal(t, tt.expectedCreateEntryIsInvoked, tt.mock.CreateEntryIsInvoked,
				"mock create entry is invoked is unexpected")
			tt.responseValidator(rw)
		})
	}
}

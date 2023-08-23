package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FluffyKebab/daytail"
	"github.com/FluffyKebab/daytail/mock"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	jwtSecretKey := "secret"

	workingPayload, err := json.Marshal(signUpRequestPayload{Name: "bob"})
	require.NoError(t, err)

	nonWorkingPayload, err := json.Marshal(struct {
		Names []string
	}{Names: []string{"bob", "amy"}})
	require.NoError(t, err)

	testCases := []struct {
		name                        string
		body                        io.Reader
		mock                        *mock.UserService
		responseValidator           func(t *testing.T, rr *httptest.ResponseRecorder)
		expectedCreateUserIsInvoked bool
		expectedStatusCode          int
	}{
		{
			name: "with valid payload",
			body: bytes.NewBuffer(workingPayload),
			mock: &mock.UserService{
				CreateUserFunc: func(d daytail.User) (int, error) {
					require.Equal(t, "bob", d.Name)
					return 1, nil
				},
			},
			responseValidator: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				var responsePayload signUpResponsePayload
				err := json.Unmarshal(rr.Body.Bytes(), &responsePayload)
				require.NoError(t, err)

				token, err := jwtauth.New("HS256", []byte(jwtSecretKey), nil).Decode(responsePayload.Token)
				require.NoError(t, err)

				value, ok := token.Get("userId")
				require.True(t, ok, "userId not present in jwt token")
				require.Equal(t, "1", fmt.Sprint(value))
			},
			expectedCreateUserIsInvoked: true,
			expectedStatusCode:          http.StatusOK,
		},
		{
			name: "with invalid payload",
			body: bytes.NewBuffer(nonWorkingPayload),
			mock: &mock.UserService{
				CreateUserFunc: func(d daytail.User) (int, error) {
					return 0, nil
				},
			},
			responseValidator: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
			},
			expectedCreateUserIsInvoked: false,
			expectedStatusCode:          http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				UserService: tt.mock,
			}
			h.InitAuth(jwtSecretKey)

			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", tt.body)
			h.signUp(rw, req)
			require.Equal(t, tt.expectedCreateUserIsInvoked, tt.mock.CrateUserIsInvoked,
				"mock create user is invoked is unexpected")
			require.Equal(t, tt.expectedStatusCode, rw.Result().StatusCode,
				"unexpected status code")
			tt.responseValidator(t, rw)
		})
	}
}

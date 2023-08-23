package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/FluffyKebab/daytail"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	h, err := initServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	go runServer(h)
	os.Exit(m.Run())
}

func TestCreateUserAndEntries(t *testing.T) {
	t.Parallel()

	b, err := json.Marshal(daytail.User{Name: "Bob"})
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/users", bytes.NewBuffer(b))
	require.NoError(t, err)
	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	var signUpResponseData struct {
		Token  string `json:"token"`
		UserId int    `json:"userId"`
	}

	err = json.NewDecoder(r.Body).Decode(&signUpResponseData)
	require.NoError(t, err)

	b, err = json.Marshal(daytail.Entry{
		Title: "my day",
		Text:  "very good",
	})
	require.NoError(t, err)

	req, err = http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:8080/users/%v/entries", signUpResponseData.UserId),
		bytes.NewBuffer(b),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "BEARER "+signUpResponseData.Token)

	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	req, err = http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:8080/users/%v/entries", signUpResponseData.UserId),
		nil,
	)

	require.NoError(t, err)
	req.Header.Set("Authorization", "BEARER "+signUpResponseData.Token)

	r, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	var entriesResponseData []struct {
		ID     int    `json:"id"`
		UserID int    `json:"userId"`
		Title  string `json:"title"`
		Text   string `json:"text"`
	}

	err = json.NewDecoder(r.Body).Decode(&entriesResponseData)
	require.NoError(t, err)
	require.Equal(t, 1, len(entriesResponseData), "entriesResponseData expected to have a length of one")
	require.Equal(t, signUpResponseData.UserId, entriesResponseData[0].UserID, "userId not equal")
	require.Equal(t, "my day", entriesResponseData[0].Title, "title wrong")
	require.Equal(t, "very good", entriesResponseData[0].Text, "text wrong")
}

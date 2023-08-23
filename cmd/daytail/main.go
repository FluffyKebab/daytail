package main

import (
	"log"

	"github.com/FluffyKebab/daytail/http"
	"github.com/FluffyKebab/daytail/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	h, err := initServer()
	if err != nil {
		return err
	}

	return runServer(h)
}

func initServer() (http.Handler, error) {
	db, err := sqlite.Migrate()
	if err != nil {
		return http.Handler{}, err
	}

	h := http.Handler{
		UserService:  sqlite.UserService{DB: db},
		EntryService: sqlite.EntryService{DB: db},
	}
	h.InitAuth("very_secret_code")

	return h, nil
}

func runServer(h http.Handler) error {
	defer h.UserService.Close()
	return http.ListenAndServe(":8080", h)
}

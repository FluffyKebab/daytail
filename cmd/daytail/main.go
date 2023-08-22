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
	defer db.Close()

	userService := sqlite.UserService{DB: db}
	entryService := sqlite.EntryService{DB: db}

	h := http.Handler{
		UserService:  userService,
		EntryService: entryService,
	}
	h.InitAuth("very_secret_code")

	return h, nil
}

func runServer(h http.Handler) error {
	return http.ListenAndServe(":8080", h)
}

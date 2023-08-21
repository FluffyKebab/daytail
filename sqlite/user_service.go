package sqlite

import (
	"database/sql"

	"github.com/FluffyKebab/daytail"
)

type UserService struct {
	*sql.DB
}

var _ daytail.UserService = UserService{}

func (s UserService) User(id int) (daytail.User, error) {
	return daytail.User{}, nil
}

func (s UserService) CreateUser(u daytail.User) (int, error) {
	var id int
	err := s.DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", u.Name).Scan(&id)
	return id, err
}

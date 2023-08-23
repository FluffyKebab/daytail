package mock

import "github.com/FluffyKebab/daytail"

type UserService struct {
	UserFunc      func(id int) (daytail.User, error)
	UserIsInvoked bool

	CreateUserFunc     func(d daytail.User) (int, error)
	CrateUserIsInvoked bool
}

var _ daytail.UserService = &UserService{}

func (s *UserService) User(id int) (daytail.User, error) {
	s.UserIsInvoked = true
	return s.UserFunc(id)
}

func (s *UserService) CreateUser(u daytail.User) (int, error) {
	s.CrateUserIsInvoked = true
	return s.CreateUserFunc(u)
}

func (s *UserService) Close() error { return nil }

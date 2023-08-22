package mock

import "github.com/FluffyKebab/daytail"

type EntryService struct {
	EntryFunc      func(id int) error
	EntryIsInvoked bool

	UserEntriesFunc      func(userId int) ([]daytail.Entry, error)
	UserEntriesIsInvoked bool

	CreateEntryFunc      func(entry daytail.Entry) (int, error)
	CreateEntryIsInvoked bool
}

var _ daytail.EntryService = &EntryService{}

func (s *EntryService) Entry(id int) error {
	s.EntryIsInvoked = true
	return s.EntryFunc(id)
}

func (s *EntryService) UserEntries(userId int) ([]daytail.Entry, error) {
	s.UserEntriesIsInvoked = true
	return s.UserEntriesFunc(userId)
}

func (s *EntryService) CreateEntry(entry daytail.Entry) (int, error) {
	s.CreateEntryIsInvoked = true
	return s.CreateEntryFunc(entry)
}

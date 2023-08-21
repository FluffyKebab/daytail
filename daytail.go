package daytail

type Entry struct {
	ID     int
	UserID int
	Title  string
	Text   string
}

type User struct {
	ID   int
	Name string
}

type UserService interface {
	User(id int) (User, error)
	CreateUser(u User) (int, error)
}

type EntryService interface {
	Entry(id int) error
	UserEntries(userId int) ([]Entry, error)
	CreateEntry(entry Entry) error
}

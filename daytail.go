package daytail

type Entry struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	User(id int) (User, error)
	CreateUser(u User) (int, error)
}

type EntryService interface {
	Entry(id int) error
	UserEntries(userId int) ([]Entry, error)
	CreateEntry(entry Entry) (int, error)
}

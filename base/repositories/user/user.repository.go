package user_repository

type User struct {
	ID                string
	Email             string
	Password          string
	Username          string
	CanonicalUsername string
}

type UserRepository interface {
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmailOrUsername(text string) (*User, error)
	CreateUser(email, password, username string) (*User, error)
	DeleteUser(userID string) error
}

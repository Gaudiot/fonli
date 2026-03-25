package user_repository

import (
	"strings"

	"github.com/google/uuid"
)

type UserRepositoryMock struct {
	Users map[string]*User
}

func (r *UserRepositoryMock) GetUserByID(id string) (*User, error) {
	user, ok := r.Users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *UserRepositoryMock) GetUserByEmail(email string) (*User, error) {
	for _, user := range r.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (r *UserRepositoryMock) GetUserByEmailOrUsername(text string) (*User, error) {
	if strings.Contains(text, "@") {
		return r.GetUserByEmail(text)
	}
	return r.GetUserByUsername(text)
}

func (r *UserRepositoryMock) GetUserByUsername(username string) (*User, error) {
	formattedUsername := formatUsername(username)
	for _, user := range r.Users {
		if user.CanonicalUsername == formattedUsername {
			return user, nil
		}
	}
	return nil, nil
}

func (r *UserRepositoryMock) CreateUser(email, password, username string) (*User, error) {
	id := uuid.New().String()
	formattedUsername := formatUsername(username)
	user := &User{
		ID:                id,
		Email:             email,
		Password:          password,
		Username:          username,
		CanonicalUsername: formattedUsername,
	}
	r.Users[id] = user
	return user, nil
}

func (r *UserRepositoryMock) DeleteUser(userID string) error {
	delete(r.Users, userID)
	return nil
}

func formatUsername(username string) string {
	formattedUsername := strings.ToLower(username)

	return formattedUsername
}

package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BCryptPasswordService struct{}

func (ps *BCryptPasswordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrGeneratingHash
	}

	return string(hash), nil
}

func (ps *BCryptPasswordService) Compare(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return ErrInvalidHash
	}
	return nil
}

package password

import "errors"

var (
	ErrGeneratingHash   = errors.New("error generating hash")
	ErrInvalidHash      = errors.New("error hashing password")
	ErrPasswordMismatch = errors.New("password mismatch")
)

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, hashedPassword string) error
}

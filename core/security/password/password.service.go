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

// MARK: - MOCK

type PasswordServiceMock struct{}

func (ps *PasswordServiceMock) Hash(password string) (string, error) {
	return password, nil
}

func (ps *PasswordServiceMock) Compare(password, hashedPassword string) error {
	if password != hashedPassword {
		return ErrPasswordMismatch
	}
	return nil
}

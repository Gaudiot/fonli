package password

import "errors"

type PasswordServiceMock struct {
	HashFunc      func(password string) (string, error)
	HashCallCount int

	CompareFunc      func(password, hashedPassword string) error
	CompareCallCount int
}

func (m *PasswordServiceMock) Hash(password string) (string, error) {
	m.HashCallCount++
	if m.HashFunc != nil {
		return m.HashFunc(password)
	}
	return "", errors.New("[Mock] not implemented")
}

func (m *PasswordServiceMock) Compare(password, hashedPassword string) error {
	m.CompareCallCount++
	if m.CompareFunc != nil {
		return m.CompareFunc(password, hashedPassword)
	}
	return errors.New("[Mock] not implemented")
}

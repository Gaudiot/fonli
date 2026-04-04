package refreshtoken_repository

import (
	"errors"
	"time"
)

type RefreshTokenRepositoryMock struct {
	CreateRefreshTokenFunc      func(token, userID string, expiresAt time.Time) (*RefreshToken, error)
	CreateRefreshTokenCallCount int

	GetByTokenFunc      func(token string) (*RefreshToken, error)
	GetByTokenCallCount int

	InvalidateUserRefreshTokensFunc      func(userID string) error
	InvalidateUserRefreshTokensCallCount int

	DeleteByTokenFunc      func(token string) error
	DeleteByTokenCallCount int

	DeleteAllByUserIDFunc      func(userID string) error
	DeleteAllByUserIDCallCount int
}

func (m *RefreshTokenRepositoryMock) CreateRefreshToken(token, userID string, expiresAt time.Time) (*RefreshToken, error) {
	m.CreateRefreshTokenCallCount++
	if m.CreateRefreshTokenFunc != nil {
		return m.CreateRefreshTokenFunc(token, userID, expiresAt)
	}
	return nil, errors.New("[Mock] not implemented")
}

func (m *RefreshTokenRepositoryMock) GetByToken(token string) (*RefreshToken, error) {
	m.GetByTokenCallCount++
	if m.GetByTokenFunc != nil {
		return m.GetByTokenFunc(token)
	}
	return nil, errors.New("[Mock] not implemented")
}

func (m *RefreshTokenRepositoryMock) InvalidateUserRefreshTokens(userID string) error {
	m.InvalidateUserRefreshTokensCallCount++
	if m.InvalidateUserRefreshTokensFunc != nil {
		return m.InvalidateUserRefreshTokensFunc(userID)
	}
	return errors.New("[Mock] not implemented")
}

func (m *RefreshTokenRepositoryMock) DeleteByToken(token string) error {
	m.DeleteByTokenCallCount++
	if m.DeleteByTokenFunc != nil {
		return m.DeleteByTokenFunc(token)
	}
	return errors.New("[Mock] not implemented")
}

func (m *RefreshTokenRepositoryMock) DeleteAllByUserID(userID string) error {
	m.DeleteAllByUserIDCallCount++
	if m.DeleteAllByUserIDFunc != nil {
		return m.DeleteAllByUserIDFunc(userID)
	}
	return errors.New("[Mock] not implemented")
}

func NewInMemoryRefreshTokenRepositoryMock(tokens map[string]*RefreshToken) *RefreshTokenRepositoryMock {
	m := &RefreshTokenRepositoryMock{}
	m.CreateRefreshTokenFunc = func(token, userID string, expiresAt time.Time) (*RefreshToken, error) {
		refreshToken := &RefreshToken{
			Token:     token,
			UserID:    userID,
			ExpiresAt: expiresAt,
			IsValid:   true,
		}
		tokens[token] = refreshToken
		return refreshToken, nil
	}
	m.GetByTokenFunc = func(token string) (*RefreshToken, error) {
		refreshToken, ok := tokens[token]
		if !ok {
			return nil, ErrInvalidRefreshToken
		}
		if refreshToken.ExpiresAt.Before(time.Now()) {
			return nil, ErrExpiredRefreshToken
		}
		return refreshToken, nil
	}
	m.InvalidateUserRefreshTokensFunc = func(userID string) error {
		for _, refreshToken := range tokens {
			if refreshToken.UserID == userID {
				refreshToken.IsValid = false
			}
		}
		return nil
	}
	m.DeleteByTokenFunc = func(token string) error {
		delete(tokens, token)
		return nil
	}
	m.DeleteAllByUserIDFunc = func(userID string) error {
		for t, refreshToken := range tokens {
			if refreshToken.UserID == userID {
				delete(tokens, t)
			}
		}
		return nil
	}
	return m
}

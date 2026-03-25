package refreshtoken_repository

import (
	"errors"
	"time"
)

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
)

type RefreshTokenRepository interface {
	CreateRefreshToken(token, userID string, expiresAt time.Time) (*RefreshToken, error)
	GetByToken(token string) (*RefreshToken, error)
	DeleteByToken(token string) error
	DeleteAllByUserID(userID string) error
}

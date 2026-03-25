package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenLifetime = 30 * time.Minute
	issuer              = "fonli"
)

var (
	ErrInvalidAccessTokenPayload = errors.New("invalid access token payload")
	ErrSigningAccessToken        = errors.New("error signing access token")
	ErrInvalidSigningMethod      = errors.New("invalid signing method")
	ErrExpiredAccessToken        = errors.New("expired access token")
	ErrInvalidAccessToken        = errors.New("invalid access token")
)

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	signingKey []byte
}

func NewTokenService(signingKey []byte) *TokenService {
	return &TokenService{
		signingKey: signingKey,
	}
}

func (tk *TokenService) GenerateAccessToken(userID string) (string, error) {
	if userID == "" {
		return "", ErrInvalidAccessTokenPayload
	}

	now := time.Now()
	claims := &accessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenLifetime)),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(tk.signingKey)
	if err != nil {
		return "", ErrSigningAccessToken
	}

	return signedToken, nil
}

func (tk *TokenService) ParseAccessToken(tokenString string) (*accessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return tk.signingKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredAccessToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidAccessToken, err)
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidAccessToken
	}

	return claims, nil
}

func (tk *TokenService) GenerateRefreshToken() string {
	refreshToken := uuid.New().String()

	return refreshToken
}

package refreshtoken_repository

import "time"

type RefreshTokenRepositoryMock struct {
	RefreshTokens map[string]*RefreshToken
}

func (r *RefreshTokenRepositoryMock) CreateRefreshToken(token, userID string, expiresAt time.Time) (*RefreshToken, error) {
	refreshToken := &RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	r.RefreshTokens[token] = refreshToken

	return refreshToken, nil
}

func (r *RefreshTokenRepositoryMock) GetByToken(token string) (*RefreshToken, error) {
	refreshToken, ok := r.RefreshTokens[token]
	if !ok {
		return nil, ErrInvalidRefreshToken
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredRefreshToken
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepositoryMock) DeleteByToken(token string) error {
	delete(r.RefreshTokens, token)

	return nil
}

func (r *RefreshTokenRepositoryMock) DeleteAllByUserID(userID string) error {
	for token, refreshToken := range r.RefreshTokens {
		if refreshToken.UserID == userID {
			delete(r.RefreshTokens, token)
		}
	}
	return nil
}

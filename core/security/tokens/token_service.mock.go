package tokens

import "errors"

type TokenServiceMock struct {
	GenerateAccessTokenFunc      func(userID string) (string, error)
	GenerateAccessTokenCallCount int

	ParseAccessTokenFunc      func(accessTokenString string) (*accessTokenClaims, error)
	ParseAccessTokenCallCount int

	GenerateRefreshTokenFunc      func() string
	GenerateRefreshTokenCallCount int
}

func (m *TokenServiceMock) GenerateAccessToken(userID string) (string, error) {
	m.GenerateAccessTokenCallCount++
	if m.GenerateAccessTokenFunc != nil {
		return m.GenerateAccessTokenFunc(userID)
	}
	return "", errors.New("[Mock] not implemented")
}

func (m *TokenServiceMock) ParseAccessToken(accessTokenString string) (*accessTokenClaims, error) {
	m.ParseAccessTokenCallCount++
	if m.ParseAccessTokenFunc != nil {
		return m.ParseAccessTokenFunc(accessTokenString)
	}
	return nil, errors.New("[Mock] not implemented")
}

func (m *TokenServiceMock) GenerateRefreshToken() string {
	m.GenerateRefreshTokenCallCount++
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc()
	}
	return ""
}

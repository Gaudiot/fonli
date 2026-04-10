package tokens

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	ParseAccessToken(accessTokenString string) (*accessTokenClaims, error)
	GenerateRefreshToken() string
}

package authentication

import (
	"strings"
	"time"

	refreshtoken_repo "gaudiot.com/fonli/base/repositories/refresh_token"
	user_repo "gaudiot.com/fonli/base/repositories/user"
	"gaudiot.com/fonli/core/security/password"
	"gaudiot.com/fonli/core/security/tokens"
	auth_validator "gaudiot.com/fonli/src/authentication/validators"
)

const (
	refreshTokenLifetime = 30 * 24 * time.Hour // 30 days
	invalidPasswordHash  = "invalidpasswordhash"
)

var (
	ErrInvalidEmail          = auth_validator.ErrInvalidEmail
	ErrPasswordTooShort      = auth_validator.ErrPasswordTooShort
	ErrPasswordTooLong       = auth_validator.ErrPasswordTooLong
	ErrPasswordMissingLetter = auth_validator.ErrPasswordMissingLetter
	ErrPasswordMissingDigit  = auth_validator.ErrPasswordMissingDigit
	ErrPasswordNoWhitespace  = auth_validator.ErrPasswordNoWhitespace

	ErrUsernameTooShort      = auth_validator.ErrUsernameTooShort
	ErrUsernameTooLong       = auth_validator.ErrUsernameTooLong
	ErrUsernameInvalidFormat = auth_validator.ErrUsernameInvalidFormat

	ErrEmailAlreadyRegistered    = auth_validator.ErrEmailAlreadyRegistered
	ErrUsernameAlreadyRegistered = auth_validator.ErrUsernameAlreadyRegistered
	ErrInvalidCredentials        = auth_validator.ErrInvalidCredentials

	ErrInvalidRefreshToken = refreshtoken_repo.ErrInvalidRefreshToken
	ErrExpiredRefreshToken = refreshtoken_repo.ErrExpiredRefreshToken
)

type AuthService struct {
	tokenService           tokens.TokenService
	passwordService        password.PasswordService
	userRepository         user_repo.UserRepository
	refreshTokenRepository refreshtoken_repo.RefreshTokenRepository
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(tokenService tokens.TokenService, passwordService password.PasswordService, userRepository user_repo.UserRepository, refreshTokenRepository refreshtoken_repo.RefreshTokenRepository) *AuthService {
	return &AuthService{
		tokenService:           tokenService,
		passwordService:        passwordService,
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (as *AuthService) generateAuthTokens(userID string) (*AuthTokens, error) {
	accessToken, err := as.tokenService.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	refreshToken := as.tokenService.GenerateRefreshToken()

	timeNow := time.Now()
	if _, err := as.refreshTokenRepository.CreateRefreshToken(refreshToken, userID, timeNow.Add(refreshTokenLifetime)); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// MARK: - SignUp
func validateSignUpInput(username, email, password string) error {
	if err := auth_validator.ValidateUsername(username); err != nil {
		return err
	}
	if err := auth_validator.ValidateEmail(email); err != nil {
		return err
	}
	if err := auth_validator.ValidatePassword(password); err != nil {
		return err
	}
	return nil
}

func (as *AuthService) SignUp(username, email, password string) (*AuthTokens, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))

	if err := validateSignUpInput(username, email, password); err != nil {
		return nil, err
	}

	if user, err := as.userRepository.GetUserByEmail(email); err != nil {
		return nil, err
	} else if user != nil {
		return nil, auth_validator.ErrEmailAlreadyRegistered
	}

	if user, err := as.userRepository.GetUserByUsername(username); err != nil {
		return nil, err
	} else if user != nil {
		return nil, auth_validator.ErrUsernameAlreadyRegistered
	}

	hashedPassword, err := as.passwordService.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := as.userRepository.CreateUser(email, hashedPassword, username)
	if err != nil {
		return nil, err
	}

	tks, err := as.generateAuthTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return tks, nil
}

// MARK: - Login
func (as *AuthService) Login(emailOrUsername, password string) (*AuthTokens, error) {
	emailOrUsername = strings.ToLower(strings.TrimSpace(emailOrUsername))

	user, err := as.userRepository.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		return nil, err
	}

	if user == nil {
		as.passwordService.Compare(password, invalidPasswordHash)
		return nil, ErrInvalidCredentials
	}

	err = as.passwordService.Compare(password, user.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	tks, err := as.generateAuthTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return tks, nil
}

// MARK: - Refresh
func (as *AuthService) Refresh(refreshToken string) (*AuthTokens, error) {
	token, err := as.refreshTokenRepository.GetByToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrInvalidRefreshToken
	}

	tks, err := as.generateAuthTokens(token.UserID)
	if err != nil {
		return nil, err
	}

	if err := as.refreshTokenRepository.DeleteByToken(refreshToken); err != nil {
		return nil, err
	}

	return tks, nil
}

// MARK: - Logout
func (as *AuthService) Logout(userID string) error {
	if err := as.refreshTokenRepository.DeleteAllByUserID(userID); err != nil {
		return err
	}

	return nil
}

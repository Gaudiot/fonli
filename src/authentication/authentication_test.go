package authentication

import (
	"errors"
	"testing"
	"time"

	refreshtoken_repository "gaudiot.com/fonli/base/repositories/refresh_token"
	repository "gaudiot.com/fonli/base/repositories/user"
	"gaudiot.com/fonli/core/security/password"
	"gaudiot.com/fonli/core/security/tokens"
)

func newPasswordServiceMock() *password.PasswordServiceMock {
	return &password.PasswordServiceMock{
		HashFunc: func(pw string) (string, error) {
			return pw, nil
		},
		CompareFunc: func(pw, hashed string) error {
			if pw != hashed {
				return password.ErrPasswordMismatch
			}
			return nil
		},
	}
}

func newTestAuthServiceEmptyUsers() *AuthService {
	emptyUsers := make(map[string]*repository.User)
	tokenServiceMock := tokens.NewTokenService([]byte("test_key"))
	userRepositoryMock := repository.UserRepositoryMock{Users: emptyUsers}
	refreshTokenRepositoryMock := refreshtoken_repository.NewInMemoryRefreshTokenRepositoryMock(make(map[string]*refreshtoken_repository.RefreshToken))
	return NewAuthService(tokenServiceMock, newPasswordServiceMock(), &userRepositoryMock, refreshTokenRepositoryMock)
}

func newTestAuthServiceWithUsers() *AuthService {
	validUsers := map[string]*repository.User{
		"id1": {
			ID:                "id1",
			Email:             "email@example.com",
			Password:          "password123",
			Username:          "John_Doe",
			CanonicalUsername: "john_doe",
		},
		"id2": {
			ID:                "id2",
			Email:             "john.doe@example.com",
			Password:          "Password123!",
			Username:          "janeDoe",
			CanonicalUsername: "janedoe",
		},
	}
	refreshTokens := map[string]*refreshtoken_repository.RefreshToken{
		"refresh_token": {
			Token:     "refresh_token",
			UserID:    "id1",
			ExpiresAt: time.Now().Add(refreshTokenLifetime),
			IsValid:   true,
		},
	}
	tokenServiceMock := tokens.NewTokenService([]byte("test_key"))
	userRepositoryMock := repository.UserRepositoryMock{Users: validUsers}
	refreshTokenRepositoryMock := refreshtoken_repository.NewInMemoryRefreshTokenRepositoryMock(refreshTokens)
	return NewAuthService(tokenServiceMock, newPasswordServiceMock(), &userRepositoryMock, refreshTokenRepositoryMock)
}

// MARK: - SignUp
func TestSignUpSuccess(t *testing.T) {
	cases := []struct {
		name     string
		username string
		email    string
		password string
	}{
		{"basic", "fulano", "email@example.com", "Password123"},
		{"all allowed symbols", "us3r_Name_1", "john.doe+test@example.com", "Password123!"},
		{"boundary inputs", "abcde", "email@example.com.br", "password____-____-____-____-____-____-____-____-____-____limit64"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceEmptyUsers()
			tokens, err := authService.SignUp(tc.username, tc.email, tc.password)
			if err != nil {
				t.Fatalf("SignUp(%q); Unexpected error: %v", tc.name, err)
			}
			if tokens == nil {
				t.Fatalf("SignUp(%q); Got nil tokens, wanted tokens not nil", tc.name)
			}
			if tokens.AccessToken == "" {
				t.Fatalf("SignUp(%q); Got empty AccessToken, wanted AccessToken not empty", tc.name)
			}
			if tokens.RefreshToken == "" {
				t.Fatalf("SignUp(%q); Got empty RefreshToken, wanted RefreshToken not empty", tc.name)
			}
		})
	}
}

func TestSignUpInvalidEmail(t *testing.T) {
	cases := []struct {
		name  string
		email string
	}{
		{"no @", "invalid_email"},
		{"no TLD", "invalid@email"},
		{"empty local part", "@invalid.email"},
		{"forbidden character", "!invalid@email.com"},
		{"whitespace", "invalid @email.com"},
	}
	username := "fulano"
	password := "Password123"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceEmptyUsers()
			_, err := authService.SignUp(username, tc.email, password)
			if err == nil {
				t.Errorf("SignUp(%q); Wanted error %v, got nil with email %q", tc.name, ErrInvalidEmail, tc.email)
			}

			if !errors.Is(err, ErrInvalidEmail) {
				t.Errorf("SignUp(%q); Wanted error %v, got %v with email %q", tc.name, ErrInvalidEmail, err, tc.email)
			}
		})
	}
}

func TestSignUpInvalidPassword(t *testing.T) {
	cases := []struct {
		name          string
		password      string
		expectedError error
	}{
		{"short", "short", ErrPasswordTooShort},
		{"long", "long____-____-____-____-____-____-____-____-____-____-____-____-password", ErrPasswordTooLong},
		{"missing letter", "123456789", ErrPasswordMissingLetter},
		{"missing digit", "NoNumbers!", ErrPasswordMissingDigit},
		{"whitespace", "Invalid password", ErrPasswordNoWhitespace},
	}
	username := "fulano"
	email := "valid@email.com"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceEmptyUsers()
			_, err := authService.SignUp(username, email, tc.password)
			if err == nil {
				t.Errorf("SignUp(%q); Wanted error %v, got nil with password %q", tc.name, tc.expectedError, tc.password)
			}

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("SignUp(%q); Wanted error %v, got %v with password %q", tc.name, tc.expectedError, err, tc.password)
			}
		})
	}
}

func TestSignUpInvalidUsername(t *testing.T) {
	cases := []struct {
		name          string
		username      string
		expectedError error
	}{
		{"too short", "abcd", ErrUsernameTooShort},
		{"too long", "abcdefghij_klmnop_qrst_uvwxyz12", ErrUsernameTooLong},
		{"whitespace", "user name1", ErrUsernameInvalidFormat},
		{"forbidden symbol", "user!name", ErrUsernameInvalidFormat},
		{"accented character", "usérnáme", ErrUsernameInvalidFormat},
	}
	email := "valid@email.com"
	password := "Password123"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceEmptyUsers()
			_, err := authService.SignUp(tc.username, email, password)
			if err == nil {
				t.Fatalf("SignUp(%q); Wanted error %v, got nil with username %q", tc.name, tc.expectedError, tc.username)
			}
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("SignUp(%q); Wanted error %v, got %v with username %q", tc.name, tc.expectedError, err, tc.username)
			}
		})
	}
}

func TestSignUpEmailAlreadyRegistered(t *testing.T) {
	cases := []struct {
		name  string
		email string
	}{
		{"exact match", "email@example.com"},
		{"different casing", "Email@Example.COM"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceWithUsers()
			_, err := authService.SignUp("fulano", tc.email, "Password123")
			if err == nil {
				t.Fatalf("SignUp(%q); Wanted error %v, got nil with email %q", tc.name, ErrEmailAlreadyRegistered, tc.email)
			}
			if !errors.Is(err, ErrEmailAlreadyRegistered) {
				t.Errorf("SignUp(%q); Wanted error %v, got %v with email %q", tc.name, ErrEmailAlreadyRegistered, err, tc.email)
			}
		})
	}
}

func TestSignUpUsernameAlreadyRegistered(t *testing.T) {
	cases := []struct {
		name     string
		username string
	}{
		{"exact match", "John_Doe"},
		{"different casing", "john_doe"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceWithUsers()
			_, err := authService.SignUp(tc.username, "new@email.com", "Password123")
			if err == nil {
				t.Fatalf("SignUp(%q); Wanted error %v, got nil with username %q", tc.name, ErrUsernameAlreadyRegistered, tc.username)
			}
			if !errors.Is(err, ErrUsernameAlreadyRegistered) {
				t.Errorf("SignUp(%q); Wanted error %v, got %v with username %q", tc.name, ErrUsernameAlreadyRegistered, err, tc.username)
			}
		})
	}
}

// MARK: - Login
func TestLoginSuccess(t *testing.T) {
	cases := []struct {
		name            string
		emailOrUsername string
		password        string
	}{
		{"with email", "email@example.com", "password123"},
		{"with email different casing", "Email@Example.COM", "password123"},
		{"with username", "John_Doe", "password123"},
		{"with username different casing", "john_doe", "password123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authService := newTestAuthServiceWithUsers()
			tokens, err := authService.Login(tc.emailOrUsername, tc.password)
			if err != nil {
				t.Fatalf("Login(%q); Unexpected error: %v", tc.emailOrUsername, err)
			}
			if tokens == nil || tokens.AccessToken == "" {
				t.Fatalf("Login(%q); Got nil or empty AccessToken, wanted non-empty", tc.name)
			}
			if tokens.RefreshToken == "" {
				t.Fatalf("Login(%q); Got empty RefreshToken, wanted non-empty", tc.name)
			}
		})
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	cases := []struct {
		name            string
		emailOrUsername string
		password        string
		expectedError   error
	}{
		{
			name:            "invalid email format",
			emailOrUsername: "invalid-email@example",
			password:        "Password123",
			expectedError:   ErrInvalidCredentials,
		},
		{
			name:            "user not found (wrong email)",
			emailOrUsername: "notfound@example.com",
			password:        "Password123",
			expectedError:   ErrInvalidCredentials,
		},
		{
			name:            "user not found (wrong username)",
			emailOrUsername: "not_found",
			password:        "Password123",
			expectedError:   ErrInvalidCredentials,
		},
		{
			name:            "incorrect password",
			emailOrUsername: "email@example.com",
			password:        "wrong_password",
			expectedError:   ErrInvalidCredentials,
		},
	}

	authService := newTestAuthServiceWithUsers()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := authService.Login(tc.emailOrUsername, tc.password)
			if err == nil {
				t.Fatalf("Login(%q); Wanted error %v, got nil", tc.name, tc.expectedError)
			}
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Login(%q); Wanted error %v, got %v", tc.name, tc.expectedError, err)
			}
		})
	}
}

// MARK: - Refresh

func TestRefreshSuccess(t *testing.T) {
	authService := newTestAuthServiceWithUsers()
	tokens, err := authService.Refresh("refresh_token")

	if err != nil {
		t.Fatalf("Refresh(); Unexpected error: %v", err)
	}
	if tokens == nil || tokens.AccessToken == "" {
		t.Fatalf("Refresh(); Got nil or empty AccessToken, wanted non-empty")
	}
	if tokens.RefreshToken == "" {
		t.Fatalf("Refresh(); Got empty RefreshToken, wanted non-empty")
	}
}

func TestRefreshRotationInvalidatesOldToken(t *testing.T) {
	authService := newTestAuthServiceWithUsers()

	newTokens, err := authService.Refresh("refresh_token")
	if err != nil {
		t.Fatalf("Refresh(); Unexpected error: %v", err)
	}

	_, err = authService.Refresh("refresh_token")
	if err == nil {
		t.Fatal("Refresh(); Wanted error reusing old refresh token, got nil")
	}
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Errorf("Refresh(); Wanted error %v, got %v", ErrInvalidRefreshToken, err)
	}

	rotatedTokens, err := authService.Refresh(newTokens.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh(); Unexpected error: %v", err)
	}
	if rotatedTokens == nil || rotatedTokens.AccessToken == "" {
		t.Fatal("Refresh(); Got nil or empty AccessToken after rotation")
	}
}

func TestRefreshInvalidRefreshToken(t *testing.T) {
	authService := newTestAuthServiceWithUsers()
	_, err := authService.Refresh("invalid_refresh_token")
	if err == nil {
		t.Fatalf("Refresh(); Wanted error %v, got nil", ErrInvalidRefreshToken)
	}
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Errorf("Refresh(); Wanted error %v, got %v", ErrInvalidRefreshToken, err)
	}
}

func TestRefreshExpiredToken(t *testing.T) {
	expiredTokens := map[string]*refreshtoken_repository.RefreshToken{
		"expired_token": {
			Token:     "expired_token",
			UserID:    "id1",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		},
	}
	tokenServiceMock := tokens.NewTokenService([]byte("test_key"))
	userRepositoryMock := repository.UserRepositoryMock{Users: make(map[string]*repository.User)}
	refreshTokenRepositoryMock := refreshtoken_repository.NewInMemoryRefreshTokenRepositoryMock(expiredTokens)
	authService := NewAuthService(tokenServiceMock, newPasswordServiceMock(), &userRepositoryMock, refreshTokenRepositoryMock)

	_, err := authService.Refresh("expired_token")
	if err == nil {
		t.Fatal("Refresh(); Wanted error for expired refresh token, got nil")
	}
	if !errors.Is(err, ErrExpiredRefreshToken) {
		t.Errorf("Refresh(); Wanted error %v, got %v", ErrExpiredRefreshToken, err)
	}
}

// MARK: - Logout

func TestLogoutSuccess(t *testing.T) {
	authService := newTestAuthServiceWithUsers()
	err := authService.Logout("id1")
	if err != nil {
		t.Fatalf("Logout(); Unexpected error: %v", err)
	}
}

func TestLogoutInvalidatesRefreshToken(t *testing.T) {
	refreshStore := map[string]*refreshtoken_repository.RefreshToken{
		"refresh_token": {
			Token:     "refresh_token",
			UserID:    "id1",
			ExpiresAt: time.Now().Add(refreshTokenLifetime),
			IsValid:   true,
		},
	}
	refreshTokensMock := refreshtoken_repository.NewInMemoryRefreshTokenRepositoryMock(refreshStore)
	tokenServiceMock := tokens.NewTokenService([]byte("test_key"))
	authService := NewAuthService(tokenServiceMock, newPasswordServiceMock(), &repository.UserRepositoryMock{Users: make(map[string]*repository.User)}, refreshTokensMock)

	err := authService.Logout("id1")
	if err != nil {
		t.Fatalf("Logout(); Unexpected error: %v", err)
	}

	rt, ok := refreshStore["refresh_token"]
	if !ok {
		t.Fatal("Logout(); Wanted refresh token entry to remain in store after invalidation")
	}
	if rt.IsValid {
		t.Errorf("Logout(); Wanted refresh token invalidated (IsValid false), got IsValid true")
	}
}

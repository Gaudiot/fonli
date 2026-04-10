package auth_validator

import (
	"regexp"
	"unicode/utf8"
)

const (
	minUsernameLength = 5
	maxUsernameLength = 30
)

// Username must start with a letter, only [a-zA-Z0-9] and "_" are allowed, no spaces or line breaks, length 5-30.
var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,28}[a-zA-Z0-9]$`)

func ValidateUsername(username string) error {
	usernameLength := utf8.RuneCountInString(username)
	if usernameLength < minUsernameLength {
		return ErrUsernameTooShort
	}
	if usernameLength > maxUsernameLength {
		return ErrUsernameTooLong
	}
	if !usernamePattern.MatchString(username) {
		return ErrUsernameInvalidFormat
	}

	return nil
}

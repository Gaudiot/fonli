package auth_validator

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

var (
	whitespacePattern = regexp.MustCompile(`\s`)
	hasLetterPattern  = regexp.MustCompile(`[a-zA-Z]`)
	hasDigitPattern   = regexp.MustCompile(`[0-9]`)
)

// ValidatePassword applies size rules, absence of whitespace, letter and digit.
func ValidatePassword(password string) error {
	passwordLength := utf8.RuneCountInString(password)
	if passwordLength < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters long (current length: %d)", ErrPasswordTooShort, passwordLength)
	}
	if passwordLength > 64 {
		return fmt.Errorf("%w: password must be at most 64 characters long (current length: %d)", ErrPasswordTooLong, passwordLength)
	}

	if whitespacePattern.MatchString(password) {
		return fmt.Errorf("%w: whitespace, tab or newline are not allowed in password", ErrPasswordNoWhitespace)
	}

	if !hasLetterPattern.MatchString(password) {
		return fmt.Errorf("%w: password must contain at least one letter (a-z or A-Z)", ErrPasswordMissingLetter)
	}
	if !hasDigitPattern.MatchString(password) {
		return fmt.Errorf("%w: password must contain at least one number (0-9)", ErrPasswordMissingDigit)
	}
	return nil
}

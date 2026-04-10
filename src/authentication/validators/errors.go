package auth_validator

import "errors"

var (
	ErrUsernameTooShort      = errors.New("username must be at least 5 characters long")
	ErrUsernameTooLong       = errors.New("username must be at most 30 characters long")
	ErrUsernameInvalidFormat = errors.New("username must contain only letters, numbers or '_'")

	ErrInvalidEmail          = errors.New("invalid email")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong       = errors.New("password must be at most 64 characters long")
	ErrPasswordMissingLetter = errors.New("password must contain at least one letter")
	ErrPasswordMissingDigit  = errors.New("password must contain at least one number")
	ErrPasswordNoWhitespace  = errors.New("password must not contain whitespace")

	ErrEmailAlreadyRegistered    = errors.New("email already registered")
	ErrUsernameAlreadyRegistered = errors.New("username already registered")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidCredentials        = errors.New("invalid credentials")
)

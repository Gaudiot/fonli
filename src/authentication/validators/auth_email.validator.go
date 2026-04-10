package auth_validator

import (
	"fmt"
	"regexp"
)

// Padrão prático: parte local (sem espaços), domínio com pelo menos um ponto e TLD final com ≥2 letras.
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._+-]*[a-zA-Z0-9])?@(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// ValidateEmail verifica o formato do email com regex.
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: email cannot be empty", ErrInvalidEmail)
	}
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("%w: the format does not correspond to a valid address (e.g. name@domain.com)", ErrInvalidEmail)
	}
	return nil
}

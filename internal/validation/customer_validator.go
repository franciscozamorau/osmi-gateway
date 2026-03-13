// internal/validation/customer_validator.go
package validation

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

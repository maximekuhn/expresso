package auth

import "net/mail"

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return EntryValidationError{
			Field:  "Email",
			Reason: err.Error(),
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	return nil
}

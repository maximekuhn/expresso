package auth

import (
	"fmt"

	"github.com/google/uuid"
)

type EntryValidationError struct {
	Field  string
	Reason string
}

func (e EntryValidationError) Error() string {
	return fmt.Sprintf(
		"EntryValidationError[field='%s', reason='%s']",
		e.Field, e.Reason,
	)
}

type EntryAlreadyExistsError struct {
	Email  string
	UserID uuid.UUID
}

func (e EntryAlreadyExistsError) Error() string {
	return fmt.Sprintf(
		"EntryAlreadyExistsError[Email=%s, UserID=%s]",
		e.Email, e.UserID,
	)
}

type EntryNotFoundError struct {
	Email string
}

func (e EntryNotFoundError) Error() string {
	return fmt.Sprintf("EntryNotFoundError[Email=%s]", e.Email)
}

type BadCredentialsError struct {
	UserID uuid.UUID
}

func (e BadCredentialsError) Error() string {
	return fmt.Sprintf("BadCredentialsError[UserID=%s]", e.UserID)
}

type SessionExpiredError struct {
	SessionID string
}

func (e SessionExpiredError) Error() string {
	return fmt.Sprintf("SessionExpiredError[SessionID=%s]", e.SessionID)
}

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

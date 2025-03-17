package user

import (
	"fmt"

	"github.com/google/uuid"
)

type UserValidationError struct {
	Field  string
	Reason string
}

func (e UserValidationError) Error() string {
	return fmt.Sprintf(
		"UserValidationError[field='%s', reason='%s']",
		e.Field, e.Reason,
	)
}

type UserAlreadyExistsError struct {
	ID uuid.UUID
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("UserAlreadyExistsError[id=%s]", e.ID)
}

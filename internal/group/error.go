package group

import (
	"fmt"

	"github.com/google/uuid"
)

type GroupValidationError struct {
	Field  string
	Reason string
}

func (e GroupValidationError) Error() string {
	return fmt.Sprintf(
		"GroupValidationError[field='%s', reason='%s']",
		e.Field, e.Reason,
	)
}

type GroupAlreadyExistsError struct {
	ID uuid.UUID
}

func (e GroupAlreadyExistsError) Error() string {
	return fmt.Sprintf("GroupAlreadyExistsError[id=%s]", e.ID)
}

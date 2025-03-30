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

type AnotherGroupWithSameNameAlreadyExistsError struct {
	Name string
}

func (e AnotherGroupWithSameNameAlreadyExistsError) Error() string {
	return fmt.Sprintf("AnotherGroupWithSameNameAlreadyExistsError[name=%s]", e.Name)
}

type AlreadyMemberOfGroupError struct {
	GroupName string
	GroupID   uuid.UUID
	UserID    uuid.UUID
	IsOwner   bool
}

func (e AlreadyMemberOfGroupError) Error() string {
	return fmt.Sprintf(
		"AlreadyMemberOfGroupError[groupname='%s', groupId:=%s, userId:=%s, isOwner=%t]",
		e.GroupName, e.GroupID, e.UserID, e.IsOwner,
	)
}

type GroupNotFoundError struct {
	GroupName string
}

func (e GroupNotFoundError) Error() string {
	return fmt.Sprintf("GroupNotFoundError[name='%s']", e.GroupName)
}

type IncorrectPasswordError struct {
	GroupID uuid.UUID
}

func (e IncorrectPasswordError) Error() string {
	return fmt.Sprintf("IncorrectPasswordError[groupId=%s]", e.GroupID)
}

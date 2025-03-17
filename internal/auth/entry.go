package auth

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	Email            string
	HashedPassword   []byte
	UserID           uuid.UUID
	SessionID        *uuid.UUID
	SessionExpiresAt *time.Time
}

func NewEntry(
	email string,
	hashedPassword []byte,
	userID uuid.UUID,
	sessionID *uuid.UUID,
	sessionExpiresAt *time.Time,
) (*Entry, error) {
	e := &Entry{
		Email:            email,
		HashedPassword:   hashedPassword,
		UserID:           userID,
		SessionID:        sessionID,
		SessionExpiresAt: sessionExpiresAt,
	}
	return e, e.validate()
}

func (e *Entry) validate() error {
	return ValidateEmail(e.Email)
}

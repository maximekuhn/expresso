package auth

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	Email            string
	HashedPassword   []byte
	UserID           uuid.UUID
	SessionID        string
	SessionExpiresAt *time.Time
}

func NewEntry(
	email string,
	hashedPassword []byte,
	userID uuid.UUID,
	sessionID string,
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
	if e.SessionExpiresAt != nil && e.SessionID == "" {
		return EntryValidationError{
			Field:  "SessionID",
			Reason: "SessionExpiresAt is defined but SessionID is empty",
		}
	}
	if e.SessionID != "" && e.SessionExpiresAt == nil {
		return EntryValidationError{
			Field:  "SessionExpiresAt",
			Reason: "SessionID is defined but SessionExpiresAt is nil",
		}
	}
	return ValidateEmail(e.Email)
}

package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/common"
)

type Service struct {
	store            Store
	datetimeProvider common.DatetimeProvider
}

func NewService(
	store Store,
	datetimeProvider common.DatetimeProvider,
) *Service {
	return &Service{
		store:            store,
		datetimeProvider: datetimeProvider,
	}
}

func (s *Service) CreateAuthEntry(ctx context.Context, email, password string, userID uuid.UUID) error {
	// TODO: validate password strength
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	e, err := NewEntry(email, hashedPassword, userID, "", nil)
	if err != nil {
		return err
	}
	return s.store.Save(ctx, *e)
}

func (s *Service) CreateSession(ctx context.Context, email, password string) (*Entry, error) {
	// retrieve entry
	entry, found, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, EntryNotFoundError{Email: email}
	}

	// compare password with hashed password
	if !checkPassword(password, entry.HashedPassword) {
		return nil, BadCredentialsError{UserID: entry.UserID}
	}

	// create session and update the entry
	sessionID := generateSessionId()
	sessionExpiresAt := s.datetimeProvider.Provide().Add(24 * time.Hour)
	newEntry, err := NewEntry(
		entry.Email,
		entry.HashedPassword,
		entry.UserID,
		sessionID,
		&sessionExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return newEntry, s.store.Update(ctx, *entry, *newEntry)
}

// IsSessionValid returns the userID and true if the session is valid.
// If the session is invalid or not found, false is returned.
// An error with more details might be returned.
func (s *Service) IsSessionValid(ctx context.Context, c *http.Cookie) (uuid.UUID, bool, error) {
	sessionId := c.Value
	entry, found, err := s.store.GetBySessionID(ctx, sessionId)
	if err != nil {
		return uuid.UUID{}, false, err
	}
	if !found {
		return uuid.UUID{}, false, nil
	}
	if !entry.IsSessionActive() {
		return uuid.UUID{}, false, nil
	}

	// session is active, check if it is valid (not expired)
	if entry.SessionExpiresAt.Before(s.datetimeProvider.Provide()) {
		return uuid.UUID{}, false, SessionExpiredError{entry.SessionID}
	}
	return entry.UserID, true, nil
}

package auth

import (
	"context"
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

func (s *Service) CreateSession(ctx context.Context, email, password string) (string, error) {
	// retrieve entry
	entry, found, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if !found {
		return "", EntryNotFoundError{Email: email}
	}

	// compare password with hashed password
	if !checkPassword(password, entry.HashedPassword) {
		return "", BadCredentialsError{UserID: entry.UserID}
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
		return "", err
	}
	return sessionID, s.store.Update(ctx, *entry, *newEntry)
}

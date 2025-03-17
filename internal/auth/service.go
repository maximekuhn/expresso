package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateAuthEntry(ctx context.Context, email, password string, userID uuid.UUID) error {
	// TODO: validate password strength
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	e, err := NewEntry(email, hashedPassword, userID, nil, nil)
	if err != nil {
		return err
	}
	return s.store.Save(ctx, *e)
}

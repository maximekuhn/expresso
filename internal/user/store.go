package user

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, u User) error
	GetById(ctx context.Context, userID uuid.UUID) (*User, bool, error)
}

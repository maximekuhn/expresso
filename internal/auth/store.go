package auth

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, e Entry) error
	Update(ctx context.Context, old, new Entry) error
	GetByEmail(ctx context.Context, email string) (*Entry, bool, error)
	GetBySessionID(ctx context.Context, id string) (*Entry, bool, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (*Entry, bool, error)
}

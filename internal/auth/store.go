package auth

import (
	"context"
)

type Store interface {
	Save(ctx context.Context, e Entry) error
	Update(ctx context.Context, old, new Entry) error
	GetByEmail(ctx context.Context, email string) (*Entry, bool, error)
}

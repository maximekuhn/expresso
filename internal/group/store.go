package group

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, g Group) error
	GetAllWhereUserIsOwner(ctx context.Context, userID uuid.UUID) ([]Group, error)
	GetAllWhereUserIsMember(ctx context.Context, userID uuid.UUID) ([]Group, error)
}

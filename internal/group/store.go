package group

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, g Group) error
	GetAllWhereUserIsOwner(ctx context.Context, userID uuid.UUID) ([]Group, error)
	GetAllWhereUserIsMember(ctx context.Context, userID uuid.UUID) ([]Group, error)
	GetByGroupName(ctx context.Context, name string) (*Group, bool, error)
	AddMember(ctx context.Context, groupID, userID uuid.UUID) error
}

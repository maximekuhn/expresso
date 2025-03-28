package group

import "context"

type Store interface {
	Save(ctx context.Context, g Group) error
}

package auth

import "context"

type Store interface {
	Save(ctx context.Context, e Entry) error
}

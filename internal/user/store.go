package user

import "context"

type Store interface {
	Save(ctx context.Context, u User) error
}

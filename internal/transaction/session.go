package transaction

import "context"

type Session interface {
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

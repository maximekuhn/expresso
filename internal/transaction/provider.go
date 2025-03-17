package transaction

import "context"

type SessionProvider interface {
	Provide(ctx context.Context) Session
}

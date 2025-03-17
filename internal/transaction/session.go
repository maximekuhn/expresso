package transaction

import "context"

type Session interface {
	// Transaction executes the function f in an atomic way. If f returns an error,
	// the transaction is rollbacked; otherwise it is committed.
	//
	// Usage example:
	//
	//  session := SessionProvider.Provide() // get a new session
	//  ctx := context.TODO()
	//  session.Transaction(ctx, func(ctx context.Context) error {
	//      // ...
	//      return nil
	//  })
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/maximekuhn/expresso/internal/transaction"
)

type sessionKey struct{}

type sqliteSession struct {
	db *sql.DB
	tx *sql.Tx
}

func newSqliteSession(db *sql.DB) *sqliteSession {
	return &sqliteSession{
		db: db,
		tx: nil,
	}
}

func sqliteSessionFromCtx(ctx context.Context, fallback *sql.DB) *sqliteSession {
	session, ok := ctx.Value(sessionKey{}).(*sqliteSession)
	if !ok {
		return newSqliteSession(fallback)
	}
	if session == nil {
		return newSqliteSession(fallback)
	}
	return session
}

func (s *sqliteSession) txStarted() bool {
	return s.tx != nil
}

func (s *sqliteSession) startTx(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	s.tx = tx
	return nil
}

func (s *sqliteSession) commit() error {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	return s.tx.Commit()
}

func (s *sqliteSession) rollback() error {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	return s.tx.Rollback()
}

func (s *sqliteSession) Transaction(ctx context.Context, f func(ctx context.Context) error) error {
	// retrieve session from context, or create a new one
	session := sqliteSessionFromCtx(ctx, s.db)

	// check if transaction already started
	// if yes, use it
	// if no, start a new one
	txAlreadyStarted := session.txStarted()
	if !txAlreadyStarted {
		if err := session.startTx(ctx); err != nil {
			return err
		}
	}

	// put transaction in the context
	transactionCtx := context.WithValue(ctx, sessionKey{}, session)

	// run f within the transaction
	if err := f(transactionCtx); err != nil {
		if rollbackErr := session.rollback(); rollbackErr != nil {
			return fmt.Errorf(
				"failed to rollback: %s, original error: %s",
				rollbackErr, err,
			)
		}
		return err
	}

	// commit only if we started transaction
	if txAlreadyStarted {
		// tx started outside this call, we don't commit here
		return nil
	}
	return session.commit()
}

func (s *sqliteSession) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if s.tx != nil {
		return s.tx.ExecContext(ctx, query, args...)
	}
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqliteSession) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if s.tx != nil {
		return s.tx.QueryRowContext(ctx, query, args...)
	}
	return s.db.QueryRowContext(ctx, query, args...)
}

type SqliteSessionProvider struct {
	db *sql.DB
}

func NewSqliteSessionProvider(db *sql.DB) *SqliteSessionProvider {
	return &SqliteSessionProvider{
		db: db,
	}
}

func (sp *SqliteSessionProvider) Provide(_ context.Context) transaction.Session {
	return newSqliteSession(sp.db)
}

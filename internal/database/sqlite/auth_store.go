package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/maximekuhn/expresso/internal/auth"
)

type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (as *AuthStore) Save(ctx context.Context, e auth.Entry) error {
	query := `
    INSERT INTO e_auth
    (user_id, email, hashed_password)
    VALUES (?, ?, ?)
    `
	res, err := sqliteSessionFromCtx(ctx, as.db).ExecContext(ctx, query, e.UserID, e.Email, e.HashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_auth.email") {
			return auth.EntryAlreadyExistsError{
				Email:  e.Email,
				UserID: e.UserID,
			}
		}
		return err
	}
	return checkRowsAffected(res, 1)
}

package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/maximekuhn/expresso/internal/user"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) Save(ctx context.Context, u user.User) error {
	query := `
    INSERT INTO e_user
    (id, name, created_at)
    VALUES (?, ?, ?)
    `
	res, err := sqliteSessionFromCtx(ctx, us.db).ExecContext(
		ctx,
		query,
		u.ID,
		u.Name,
		u.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_user.id") {
			return user.UserAlreadyExistsError{ID: u.ID}
		}
		return err
	}
	return checkRowsAffected(res, 1)
}

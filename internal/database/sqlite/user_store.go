package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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
func (us *UserStore) GetById(ctx context.Context, userID uuid.UUID) (*user.User, bool, error) {
	query := `
    SELECT name, created_at
    FROM e_user
    WHERE id = ?
    `

	row := sqliteSessionFromCtx(ctx, us.db).QueryRowContext(ctx, query, userID)

	var name string
	var createdAt time.Time

	if err := row.Scan(&name, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	u, err := user.New(userID, name, createdAt)
	return u, true, err
}

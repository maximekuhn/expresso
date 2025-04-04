package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (as *AuthStore) Update(ctx context.Context, old, new auth.Entry) error {
	if new.UserID != old.UserID {
		return errors.New("changing UserID is not supported/forbidden as it is the primary key")
	}

	updates := make([]string, 0)
	args := make([]interface{}, 0)

	if old.Email != new.Email {
		updates = append(updates, "email = ? ")
		args = append(args, new.Email)
	}

	if !slices.Equal(old.HashedPassword, new.HashedPassword) {
		updates = append(updates, "hashed_password = ? ")
		args = append(args, new.HashedPassword)
	}

	if old.SessionID != new.SessionID {
		updates = append(updates, "session_id = ? ")
		if new.SessionID != "" {
			args = append(args, new.SessionID)
		} else {
			args = append(args, nil)
		}
	}

	if old.SessionExpiresAt != new.SessionExpiresAt {
		updates = append(updates, "session_expires_at = ? ")
		if new.SessionExpiresAt != nil {
			args = append(args, new.SessionExpiresAt.UTC())
		} else {
			args = append(args, nil)
		}
	}

	if len(updates) == 0 {
		// no update to perform
		return nil
	}

	query := fmt.Sprintf("update e_auth set %s where user_id = ?", strings.Join(updates, ", "))
	args = append(args, new.UserID)
	res, err := sqliteSessionFromCtx(ctx, as.db).ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return checkRowsAffected(res, 1)
}
func (as *AuthStore) GetByEmail(ctx context.Context, email string) (*auth.Entry, bool, error) {
	query := `
	   SELECT
	   user_id, hashed_password, session_id, session_expires_at
	   FROM e_auth
	   WHERE email = ?
	   `

	row := sqliteSessionFromCtx(ctx, as.db).QueryRowContext(ctx, query, email)

	var userId uuid.UUID
	var hashedPassword []byte
	var sessionId sql.NullString
	var sessionExpiresAt sql.NullTime

	if err := row.Scan(&userId, &hashedPassword, &sessionId, &sessionExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	sId := ""
	if sessionId.Valid {
		sId = sessionId.String
	}

	var sExpiresAt *time.Time
	if sessionExpiresAt.Valid {
		sExpiresAt = &sessionExpiresAt.Time
	}

	e, err := auth.NewEntry(email, hashedPassword, userId, sId, sExpiresAt)
	if err != nil {
		return nil, false, DataCorruptedError{
			Type:     "auth.Entry",
			Original: err,
		}
	}
	return e, true, nil
}

func (as *AuthStore) GetBySessionID(ctx context.Context, id string) (*auth.Entry, bool, error) {
	query := `
	   SELECT
	   user_id, hashed_password, email, session_expires_at
	   FROM e_auth
	   WHERE session_id = ?
	   `

	row := sqliteSessionFromCtx(ctx, as.db).QueryRowContext(ctx, query, id)

	var userId uuid.UUID
	var hashedPassword []byte
	var email string
	var sessionExpiresAt sql.NullTime

	if err := row.Scan(&userId, &hashedPassword, &email, &sessionExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var sExpiresAt *time.Time
	if sessionExpiresAt.Valid {
		sExpiresAt = &sessionExpiresAt.Time
	}

	e, err := auth.NewEntry(email, hashedPassword, userId, id, sExpiresAt)
	if err != nil {
		return nil, false, DataCorruptedError{
			Type:     "auth.Entry",
			Original: err,
		}
	}
	return e, true, nil
}

func (as *AuthStore) GetByUserID(ctx context.Context, id uuid.UUID) (*auth.Entry, bool, error) {
	query := `
	   SELECT
	   session_id, hashed_password, email, session_expires_at
	   FROM e_auth
	   WHERE user_id = ?
	   `

	row := sqliteSessionFromCtx(ctx, as.db).QueryRowContext(ctx, query, id)

	var sessionId sql.NullString
	var hashedPassword []byte
	var email string
	var sessionExpiresAt sql.NullTime

	if err := row.Scan(&sessionId, &hashedPassword, &email, &sessionExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var sExpiresAt *time.Time
	if sessionExpiresAt.Valid {
		sExpiresAt = &sessionExpiresAt.Time
	}

	sId := ""
	if sessionId.Valid {
		sId = sessionId.String
	}

	e, err := auth.NewEntry(email, hashedPassword, id, sId, sExpiresAt)
	if err != nil {
		return nil, false, DataCorruptedError{
			Type:     "auth.Entry",
			Original: err,
		}
	}
	return e, true, nil
}

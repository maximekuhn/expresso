package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/maximekuhn/expresso/internal/group"
)

type GroupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) *GroupStore {
	return &GroupStore{
		db: db,
	}
}

func (gs *GroupStore) Save(ctx context.Context, g group.Group) error {
	// this function should be used only when creating a group, meaning no
	// members can be present yet
	// This breaks a little bit the abstraction, but we don't want to insert members here
	if len(g.Members) != 0 {
		return fmt.Errorf("expected group to have 0 member, but found %d", len(g.Members))
	}
	query := `
    INSERT INTO e_group
    (id, name, owner_id, created_at, hashed_password)
    VALUES (?, ?, ?, ?, ?)
    `
	res, err := sqliteSessionFromCtx(ctx, gs.db).ExecContext(ctx, query, g.ID, g.Name, g.Owner, g.CreatedAt, g.HashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_group.id") {
			return group.GroupAlreadyExistsError{ID: g.ID}
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_group.name") {
			return group.AnotherGroupWithSameNameAlreadyExistsError{Name: g.Name}
		}
		return err
	}
	return checkRowsAffected(res, 1)
}

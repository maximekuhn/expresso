package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
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
	dbExecutor := sqliteSessionFromCtx(ctx, gs.db)

	query := `
    INSERT INTO e_group
    (id, name, owner_id, created_at, hashed_password)
    VALUES (?, ?, ?, ?, ?)
    `

	res, err := dbExecutor.ExecContext(ctx, query, g.ID, g.Name, g.Owner, g.CreatedAt, g.HashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_group.id") {
			return group.GroupAlreadyExistsError{ID: g.ID}
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: e_group.name") {
			return group.AnotherGroupWithSameNameAlreadyExistsError{Name: g.Name}
		}
		return err
	}

	if len(g.Members) > 0 {
		memberQuery := `INSERT INTO e_group_member (group_id, user_id) VALUES `

		args := make([]interface{}, 0, len(g.Members)*2)
		placeholders := make([]string, 0, len(g.Members))
		for _, userID := range g.Members {
			placeholders = append(placeholders, "(?, ?)")
			args = append(args, g.ID, userID)
		}

		memberQuery += strings.Join(placeholders, ", ")
		_, err = dbExecutor.ExecContext(ctx, memberQuery, args...)
		if err != nil {
			return err
		}
	}

	// TODO: same for members query
	return checkRowsAffected(res, 1)
}

func (gs *GroupStore) GetAllWhereUserIsOwner(ctx context.Context, userID uuid.UUID) ([]group.Group, error) {
	query := `
    SELECT id, name, created_at, hashed_password
    FROM e_group
    WHERE owner_id = ?
    `

	rows, err := sqliteSessionFromCtx(ctx, gs.db).db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]group.Group, 0)
	for rows.Next() {
		var id uuid.UUID
		var name string
		var createdAt time.Time
		var hashedPassword []byte

		if err := rows.Scan(&id, &name, &createdAt, &hashedPassword); err != nil {
			return nil, err
		}

		// TODO: fetch members
		g, err := group.New(
			id,
			name,
			userID,
			make([]uuid.UUID, 0),
			hashedPassword,
			createdAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *g)
	}
	return groups, nil
}

func (gs *GroupStore) GetAllWhereUserIsMember(ctx context.Context, userID uuid.UUID) ([]group.Group, error) {
	query := `
    SELECT g.id, g.owner_id, g.name, g.created_at, g.hashed_password
    FROM e_group g
    LEFT JOIN e_group_member gm ON g.id = gm.group_id
    WHERE gm.user_id = ?
    `

	rows, err := sqliteSessionFromCtx(ctx, gs.db).db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]group.Group, 0)
	for rows.Next() {
		var id uuid.UUID
		var ownerId uuid.UUID
		var name string
		var createdAt time.Time
		var hashedPassword []byte

		if err := rows.Scan(&id, &ownerId, &name, &createdAt, &hashedPassword); err != nil {
			return nil, err
		}

		// TODO: fetch other members
		members := make([]uuid.UUID, 0)
		members = append(members, userID)

		g, err := group.New(id, name, ownerId, members, hashedPassword, createdAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *g)
	}
	return groups, nil
}

func (gs *GroupStore) GetByGroupName(ctx context.Context, name string) (*group.Group, bool, error) {
	query := `
    SELECT g.id, g.owner_id, g.created_at, g.hashed_password, gm.user_id
    FROM e_group g
    LEFT JOIN e_group_member gm ON g.id = gm.group_id
    WHERE g.name = ?
    `

	rows, err := sqliteSessionFromCtx(ctx, gs.db).db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var g *group.Group
	for rows.Next() {
		var id uuid.UUID
		var ownerId uuid.UUID
		var createdAt time.Time
		var hashedPassword []byte
		var memberId uuid.UUID

		if err := rows.Scan(&id, &ownerId, &createdAt, &hashedPassword, &memberId); err != nil {
			return nil, false, err
		}

		if g == nil {
			members := make([]uuid.UUID, 0)
			members = append(members, memberId)
			gr, err := group.New(id, name, ownerId, members, hashedPassword, createdAt)
			if err != nil {
				return nil, false, err
			}
			g = gr
			continue
		}

		g.Members = append(g.Members, memberId)
	}

	if g == nil {
		return nil, false, nil
	}
	return g, true, nil
}

func (gs *GroupStore) AddMember(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `
    INSERT INTO e_group_member (user_id, group_id)
    VALUES (?, ?)
    `
	res, err := sqliteSessionFromCtx(ctx, gs.db).db.ExecContext(ctx, query, userID, groupID)
	if err != nil {
		return err
	}
	return checkRowsAffected(res, 1)
}

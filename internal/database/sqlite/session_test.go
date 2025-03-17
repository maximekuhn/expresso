package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionOk(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	session := NewSqliteSessionProvider(db).Provide()
	ctx := context.TODO()
	txErr := session.Transaction(ctx, func(ctx context.Context) error {
		if err := addRowInMigrationTableAndReturnOk(ctx, db, 10_000); err != nil {
			return err
		}
		if err := addRowInMigrationTableAndReturnOk(ctx, db, 10_001); err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, txErr)
}

func TestSessionError(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	session := NewSqliteSessionProvider(db).Provide()
	ctx := context.TODO()
	txErr := session.Transaction(ctx, func(ctx context.Context) error {
		if err := addRowInMigrationTableAndReturnOk(ctx, db, 10_000); err != nil {
			return err
		}
		return errors.New("oh no, something failed :(")
	})
	assert.Error(t, txErr, "oh no, something failed :(")

	// check that the data was not inserted
	rows, err := db.Query("select version from e_migration")
	assert.NoError(t, err)
	for rows.Next() {
		var version int
		assert.NoError(t, rows.Scan(&version))
		assert.NotEqual(t, 10_000, version)
	}
}

func TestSessionError2(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	session := NewSqliteSessionProvider(db).Provide()
	ctx := context.TODO()
	txErr := session.Transaction(ctx, func(ctx context.Context) error {
		if err := addRowInMigrationTableAndReturnOk(ctx, db, 10_000); err != nil {
			return err
		}

		// error in the middle of 2 operations
		if err := returnError(); err != nil {
			return err
		}

		if err := addRowInMigrationTableAndReturnOk(ctx, db, 10_001); err != nil {
			return err
		}

		return nil
	})
	assert.Error(t, txErr, returnError().Error())

	// check that the data was not inserted
	rows, err := db.Query("select version from e_migration")
	assert.NoError(t, err)
	for rows.Next() {
		var version int
		assert.NoError(t, rows.Scan(&version))
		assert.NotEqual(t, 10_000, version)
		assert.NotEqual(t, 10_001, version)
	}
}

func addRowInMigrationTableAndReturnOk(ctx context.Context, db *sql.DB, version int) error {
	session := sqliteSessionFromCtx(ctx, db)

	query := "insert into e_migration (version, applied_at) values (?, ?)"
	res, err := session.ExecContext(ctx, query, version, time.Now().UTC())
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d row(s)", c)
	}
	return nil
}

func returnError() error {
	return errors.New("something went wrong")
}

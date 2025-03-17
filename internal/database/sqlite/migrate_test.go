package sqlite

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	db := createTmpDb()
	defer db.Close()

	assert.NoError(t, Migrate(db))
	now := time.Now().UTC()

	query := `select version, applied_at from e_migration`
	rows, err := db.Query(query)
	assert.NoError(t, err)
	assert.True(t, rows.Next())

	var version int
	var appliedAt time.Time
	assert.NoError(t, rows.Scan(&version, &appliedAt))
	assert.True(t, version >= 1, "version should be at least 1")
	assert.True(
		t, time.Now().UTC().Sub(now).Seconds() < 1,
		"applied_at should be less than a second ago",
	)
}

func createTmpDb() *sql.DB {
	f, err := os.CreateTemp("", "test-db-*.sqlite3")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		panic(err)
	}
	return db
}

func createTmpDbWithAllMigrationsApplied() *sql.DB {
	db := createTmpDb()
	if err := Migrate(db); err != nil {
		db.Close()
		panic(err)
	}
	return db
}

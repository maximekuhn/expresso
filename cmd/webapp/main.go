package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/expresso/internal/database/sqlite"
	"github.com/maximekuhn/expresso/internal/webapp"
)

func main() {
	db := setupDB()
	defer db.Close()
	webapp.Run(db)
}

func setupDB() *sql.DB {
	if len(os.Args) < 2 {
		log.Fatal("Database file path is required as the first argument.")
	}
	dbFile := os.Args[1]
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open SQLite3 database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := sqlite.Migrate(db); err != nil {
		db.Close()
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	return db
}

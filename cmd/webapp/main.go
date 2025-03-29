package main

import (
	"database/sql"
	"io"
	"log"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/expresso/internal/database/sqlite"
	"github.com/maximekuhn/expresso/internal/webapp"
)

func main() {
	db := setupDB()
	defer db.Close()

	// TODO: should be configurable
	logFile, err := os.OpenFile("./e2e/logs-test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	logsOutput := io.MultiWriter(os.Stdout, logFile)
	l := slog.New(slog.NewJSONHandler(logsOutput, nil))

	// TODO: handle prod deployment
	isProd := false
	webapp.Run(db, l, isProd)

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

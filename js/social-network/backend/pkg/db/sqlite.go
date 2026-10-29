package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migrateSqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./social.db"
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	runMigrations(db)

	return db
}

func FixInvitedMembers(db *sql.DB) {
	_, err := db.Exec(`UPDATE group_members SET status = 'accepted' WHERE status = 'invited'`)
	if err != nil {
		log.Printf("FixInvitedMembers: %v", err)
	}
}

func runMigrations(db *sql.DB) {
	driver, err := migrateSqlite3.WithInstance(db, &migrateSqlite3.Config{})
	if err != nil {
		log.Fatal("migrate driver:", err)
	}

	// Try Docker path first, fall back to local dev path
	migrationPath := "file://pkg/db/migrations/sqlite"
	if _, statErr := os.Stat("/app/pkg/db/migrations/sqlite"); statErr == nil {
		migrationPath = "file:///app/pkg/db/migrations/sqlite"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("migrate init:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate up:", err)
	}
}

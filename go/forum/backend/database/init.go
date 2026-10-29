package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    *sql.DB
}

func InitDB() (*DB, error) {
    db, err := sql.Open("sqlite3", "./ForoumDatabase.db")
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    database := &DB{db}
    if err := database.createTables(); err != nil {
        return nil, err
    }

    if err := database.SeedCategories(); err != nil {
        return nil, err
    }
    
    return database, nil
}

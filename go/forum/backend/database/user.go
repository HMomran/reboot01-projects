package database

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt"
)

func (db *DB) CheckEmailExists(email string) (bool, error) {
    var existing string
    err := db.QueryRow("SELECT email FROM CLIENT WHERE email = ?", email).Scan(&existing)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

func (db *DB) CreateUser(name, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    _, err = db.Exec("INSERT INTO CLIENT (name, email, password) VALUES (?, ?, ?)", name, email, hashedPassword)
    return err
}

func (db *DB) GetUserByEmail(email string) (int, string, error) {
    var id int
    var hashedPassword string
    err := db.QueryRow("SELECT client_id, password FROM CLIENT WHERE email = ?", email).Scan(&id, &hashedPassword)
    return id, hashedPassword, err
}

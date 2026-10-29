package database

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

func (db *DB) GetUserIDFromSession(sessionID string) (int, error) {
    var userID int
    var expiry time.Time
    
    err := db.QueryRow("SELECT client_id, expires_at FROM SESSION WHERE session_id = ?", 
        sessionID).Scan(&userID, &expiry)
    if err != nil {
        return 0, err
    }
    
    if time.Now().After(expiry) {
        return 0, errors.New("session expired")
    }
    
    return userID, nil
}

func (db *DB) DeleteSession(sessionID string) error {
    _, err := db.Exec("DELETE FROM SESSION WHERE session_id = ?", sessionID)
    return err
}

func (db *DB) DeleteUserSessions(userID int) error {
    _, err := db.Exec("DELETE FROM SESSION WHERE client_id = ?", userID)
    return err
}

func (db *DB) CreateSession(userID int, duration time.Duration) (string, time.Time, error) {
    sessionID := uuid.New().String()
    expiry := time.Now().Add(duration)
    
    _, err := db.Exec("INSERT INTO SESSION (session_id, client_id, expires_at) VALUES (?, ?, ?)", 
                     sessionID, userID, expiry)
    if err != nil {
        return "", time.Time{}, err
    }
    
    return sessionID, expiry, nil
}

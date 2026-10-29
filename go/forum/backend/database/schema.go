package database

func (db *DB) createTables() error {
    _, err := db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {
        return err
    }

    tables := []string{
        `CREATE TABLE IF NOT EXISTS CLIENT (
            client_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );`,
        
        `CREATE TABLE IF NOT EXISTS CATEGORY (
            category_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            description TEXT
        );`,
        
        `CREATE TABLE IF NOT EXISTS POST (
            post_id INTEGER PRIMARY KEY AUTOINCREMENT,
            client_id INTEGER,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (client_id) REFERENCES CLIENT(client_id)
        );`,
        
        `CREATE TABLE IF NOT EXISTS POST_CATEGORY (
            post_id INTEGER,
            category_id INTEGER,
            PRIMARY KEY (post_id, category_id),
            FOREIGN KEY (post_id) REFERENCES POST(post_id),
            FOREIGN KEY (category_id) REFERENCES CATEGORY(category_id)
        );`,
        
        `CREATE TABLE IF NOT EXISTS POST_REACTION (
            post_id INTEGER,
            client_id INTEGER,
            is_like BOOLEAN NOT NULL,
            PRIMARY KEY (post_id, client_id),
            FOREIGN KEY (post_id) REFERENCES POST(post_id),
            FOREIGN KEY (client_id) REFERENCES CLIENT(client_id)
        );`,
        
        `CREATE TABLE IF NOT EXISTS COMMENT (
            comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            client_id INTEGER,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES POST(post_id),
            FOREIGN KEY (client_id) REFERENCES CLIENT(client_id)
        );`,
        
        `CREATE TABLE IF NOT EXISTS COMMENT_REACTION (
            reaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
            comment_id INTEGER NOT NULL,
            client_id INTEGER NOT NULL,
            is_like BOOLEAN NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (comment_id) REFERENCES COMMENT(comment_id) ON DELETE CASCADE,
            FOREIGN KEY (client_id) REFERENCES CLIENT(client_id) ON DELETE CASCADE,
            UNIQUE(comment_id, client_id)
        );`,
        
        `CREATE TABLE IF NOT EXISTS SESSION (
            session_id TEXT PRIMARY KEY,
            client_id INTEGER,
            expires_at DATETIME NOT NULL,
            FOREIGN KEY (client_id) REFERENCES CLIENT(client_id)
        );`,
    }

    for _, table := range tables {
        if _, err := db.Exec(table); err != nil {
            return err
        }
    }

    return nil
}

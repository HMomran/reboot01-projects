CREATE TABLE IF NOT EXISTS group_comments (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES group_posts(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
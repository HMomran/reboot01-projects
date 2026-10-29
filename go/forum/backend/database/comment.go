package database

func (db *DB) CreateComment(postID, clientID int, content string) error {
    _, err := db.Exec(`
        INSERT INTO COMMENT (post_id, client_id, content)
        VALUES (?, ?, ?)`,
        postID, clientID, content)
    return err
}

func (db *DB) GetCommentsByPost(postID int) ([]Comment, error) {
    query := `
    SELECT 
        c.comment_id, 
        c.post_id, 
        c.client_id, 
        c.content, 
        cl.name as author, 
        c.created_at,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM COMMENT c
    JOIN CLIENT cl ON c.client_id = cl.client_id
    LEFT JOIN (
        SELECT comment_id, COUNT(*) as count 
        FROM COMMENT_REACTION 
        WHERE is_like = 1 
        GROUP BY comment_id
    ) likes ON c.comment_id = likes.comment_id
    LEFT JOIN (
        SELECT comment_id, COUNT(*) as count 
        FROM COMMENT_REACTION 
        WHERE is_like = 0 
        GROUP BY comment_id
    ) dislikes ON c.comment_id = dislikes.comment_id
    WHERE c.post_id = ?
    ORDER BY c.created_at ASC`
    
    rows, err := db.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []Comment
    for rows.Next() {
        var comment Comment
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.ClientID, 
                        &comment.Content, &comment.Author, &comment.CreatedAt,
                        &comment.Likes, &comment.Dislikes)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }
    return comments, nil
}

package database

func (db *DB) AddReaction(postID, clientID int, isLike bool) error {
    var existingIsLike bool
    err := db.QueryRow("SELECT is_like FROM POST_REACTION WHERE post_id = ? AND client_id = ?", 
                      postID, clientID).Scan(&existingIsLike)
    
    if err == nil {
        if existingIsLike == isLike {
            _, err = db.Exec("DELETE FROM POST_REACTION WHERE post_id = ? AND client_id = ?", 
                           postID, clientID)
            return err
        } else {
            _, err = db.Exec("UPDATE POST_REACTION SET is_like = ? WHERE post_id = ? AND client_id = ?", 
                           isLike, postID, clientID)
            return err
        }
    }
    
    _, err = db.Exec("INSERT INTO POST_REACTION (post_id, client_id, is_like) VALUES (?, ?, ?)", 
                    postID, clientID, isLike)
    return err
}

func (db *DB) AddCommentReaction(commentID, userID int, isLike bool) error {
    var existingIsLike bool
    err := db.QueryRow("SELECT is_like FROM COMMENT_REACTION WHERE comment_id = ? AND client_id = ?", 
                      commentID, userID).Scan(&existingIsLike)
    
    if err == nil {
        if existingIsLike == isLike {
            _, err = db.Exec("DELETE FROM COMMENT_REACTION WHERE comment_id = ? AND client_id = ?", 
                           commentID, userID)
            return err
        } else {
            _, err = db.Exec("UPDATE COMMENT_REACTION SET is_like = ? WHERE comment_id = ? AND client_id = ?", 
                           isLike, commentID, userID)
            return err
        }
    }
    
    _, err = db.Exec("INSERT INTO COMMENT_REACTION (comment_id, client_id, is_like) VALUES (?, ?, ?)", 
                    commentID, userID, isLike)
    return err
}

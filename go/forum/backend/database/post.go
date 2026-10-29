package database

import "fmt"

func (db *DB) GetAllPosts() ([]Post, error) {
    query := `
    SELECT 
        p.post_id, 
        p.client_id, 
        p.title, 
        p.content, 
        p.created_at,
        c.name,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM POST p
    JOIN CLIENT c ON p.client_id = c.client_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 1 
        GROUP BY post_id
    ) likes ON p.post_id = likes.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 0 
        GROUP BY post_id
    ) dislikes ON p.post_id = dislikes.post_id
    ORDER BY p.created_at DESC`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.ClientID, &post.Title, &post.Content, 
            &post.CreatedAt, &post.Author, &post.Likes, &post.Dislikes)
        if err != nil {
            return nil, err
        }
        
        categories, err := db.GetPostCategories(post.ID)
        if err != nil {
            return nil, err
        }
        post.Categories = categories
        
        posts = append(posts, post)
    }
    return posts, nil
}

func (db *DB) CreatePost(clientID int, title, content string) error {
    _, err := db.Exec("INSERT INTO POST (client_id, title, content) VALUES (?, ?, ?)", 
        clientID, title, content)
    return err
}

func (db *DB) CreatePostWithCategories(clientID int, title, content string, categoryIDs []int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    result, err := tx.Exec("INSERT INTO POST (client_id, title, content) VALUES (?, ?, ?)", 
        clientID, title, content)
    if err != nil {
        return err
    }

    postID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    for _, categoryID := range categoryIDs {
        _, err = tx.Exec("INSERT INTO POST_CATEGORY (post_id, category_id) VALUES (?, ?)", 
            postID, categoryID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (db *DB) GetPostsByCategory(categoryID int) ([]Post, error) {
    query := `
    SELECT DISTINCT
        p.post_id, 
        p.client_id, 
        p.title, 
        p.content, 
        p.created_at,
        c.name,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM POST p
    JOIN CLIENT c ON p.client_id = c.client_id
    JOIN POST_CATEGORY pc ON p.post_id = pc.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 1 
        GROUP BY post_id
    ) likes ON p.post_id = likes.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 0 
        GROUP BY post_id
    ) dislikes ON p.post_id = dislikes.post_id
    WHERE pc.category_id = ?
    ORDER BY p.created_at DESC`
    
    rows, err := db.Query(query, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.ClientID, &post.Title, &post.Content, 
            &post.CreatedAt, &post.Author, &post.Likes, &post.Dislikes)
        if err != nil {
            return nil, err
        }
        
        categories, err := db.GetPostCategories(post.ID)
        if err != nil {
            return nil, err
        }
        post.Categories = categories
        
        posts = append(posts, post)
    }
    return posts, nil
}

func (db *DB) GetPostsByMultipleCategories(categoryIDs []int) ([]Post, error) {
    if len(categoryIDs) == 0 {
        return db.GetAllPosts()
    }
    
    placeholders := ""
    for i := range categoryIDs {
        if i > 0 {
            placeholders += ","
        }
        placeholders += "?"
    }
    
    query := fmt.Sprintf(`
    SELECT DISTINCT
        p.post_id, 
        p.client_id, 
        p.title, 
        p.content, 
        p.created_at,
        c.name,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM POST p
    JOIN CLIENT c ON p.client_id = c.client_id
    JOIN POST_CATEGORY pc ON p.post_id = pc.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 1 
        GROUP BY post_id
    ) likes ON p.post_id = likes.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 0 
        GROUP BY post_id
    ) dislikes ON p.post_id = dislikes.post_id
    WHERE pc.category_id IN (%s)
    ORDER BY p.created_at DESC`, placeholders)
    
    args := make([]interface{}, len(categoryIDs))
    for i, id := range categoryIDs {
        args[i] = id
    }
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.ClientID, &post.Title, &post.Content, 
            &post.CreatedAt, &post.Author, &post.Likes, &post.Dislikes)
        if err != nil {
            return nil, err
        }
        
        categories, err := db.GetPostCategories(post.ID)
        if err != nil {
            return nil, err
        }
        post.Categories = categories
        
        posts = append(posts, post)
    }
    return posts, nil
}

func (db *DB) GetPostsByUser(userID int) ([]Post, error) {
    query := `
    SELECT 
        p.post_id, 
        p.client_id, 
        p.title, 
        p.content, 
        p.created_at,
        c.name,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM POST p
    JOIN CLIENT c ON p.client_id = c.client_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 1 
        GROUP BY post_id
    ) likes ON p.post_id = likes.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 0 
        GROUP BY post_id
    ) dislikes ON p.post_id = dislikes.post_id
    WHERE p.client_id = ?
    ORDER BY p.created_at DESC`
    
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.ClientID, &post.Title, &post.Content, 
            &post.CreatedAt, &post.Author, &post.Likes, &post.Dislikes)
        if err != nil {
            return nil, err
        }
        
        categories, err := db.GetPostCategories(post.ID)
        if err != nil {
            return nil, err
        }
        post.Categories = categories
        
        posts = append(posts, post)
    }
    return posts, nil
}

func (db *DB) GetLikedPostsByUser(userID int) ([]Post, error) {
    query := `
    SELECT DISTINCT
        p.post_id, 
        p.client_id, 
        p.title, 
        p.content, 
        p.created_at,
        c.name,
        COALESCE(likes.count, 0) as likes,
        COALESCE(dislikes.count, 0) as dislikes
    FROM POST p
    JOIN CLIENT c ON p.client_id = c.client_id
    JOIN POST_REACTION pr ON p.post_id = pr.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 1 
        GROUP BY post_id
    ) likes ON p.post_id = likes.post_id
    LEFT JOIN (
        SELECT post_id, COUNT(*) as count 
        FROM POST_REACTION 
        WHERE is_like = 0 
        GROUP BY post_id
    ) dislikes ON p.post_id = dislikes.post_id
    WHERE pr.client_id = ? AND pr.is_like = 1
    ORDER BY p.created_at DESC`
    
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.ClientID, &post.Title, &post.Content, 
            &post.CreatedAt, &post.Author, &post.Likes, &post.Dislikes)
        if err != nil {
            return nil, err
        }
        
        categories, err := db.GetPostCategories(post.ID)
        if err != nil {
            return nil, err
        }
        post.Categories = categories
        
        posts = append(posts, post)
    }
    return posts, nil
}

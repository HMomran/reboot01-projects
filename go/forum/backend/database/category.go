package database

func (db *DB) GetAllCategories() ([]Category, error) {
    query := `SELECT category_id, name, description FROM CATEGORY ORDER BY name`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var cat Category
        err := rows.Scan(&cat.ID, &cat.Name, &cat.Description)
        if err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    return categories, nil
}

func (db *DB) CreateCategory(name, description string) error {
    _, err := db.Exec("INSERT INTO CATEGORY (name, description) VALUES (?, ?)", name, description)
    return err
}

func (db *DB) SeedCategories() error {
    categories := []struct {
        name, description string
    }{
        {"General", "General discussions and topics"},
        {"Technology", "Tech news, programming, and gadgets"},
        {"Gaming", "Video games and gaming discussions"},
        {"Sports", "Sports news and discussions"},
        {"Movies", "Movies, TV shows, and entertainment"},
        {"Music", "Music discussions and recommendations"},
        {"Books", "Book reviews and reading discussions"},
        {"Science", "Scientific discoveries and discussions"},
    }

    for _, cat := range categories {
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM CATEGORY WHERE name = ?", cat.name).Scan(&count)
        if err != nil {
            return err
        }
        
        if count == 0 {
            err = db.CreateCategory(cat.name, cat.description)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (db *DB) GetPostCategories(postID int) ([]string, error) {
    query := `
    SELECT c.name 
    FROM CATEGORY c 
    JOIN POST_CATEGORY pc ON c.category_id = pc.category_id 
    WHERE pc.post_id = ?`
    
    rows, err := db.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var name string
        err := rows.Scan(&name)
        if err != nil {
            return nil, err
        }
        categories = append(categories, name)
    }
    return categories, nil
}

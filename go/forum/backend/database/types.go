package database

type Category struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type Post struct {
    ID         int      `json:"id"`
    ClientID   int      `json:"client_id"`
    Title      string   `json:"title"`
    Content    string   `json:"content"`
    CreatedAt  string   `json:"created_at"`
    Author     string   `json:"author"`
    Likes      int      `json:"likes"`
    Dislikes   int      `json:"dislikes"`
    Categories []string `json:"categories"`
    Comments    []Comment
}

type Comment struct {
    ID        int
    PostID    int
    ClientID  int
    Content   string
    Author    string
    CreatedAt string
    Likes     int
    Dislikes  int
}

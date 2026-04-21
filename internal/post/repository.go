package post

import "database/sql"

type PostRepository struct {
	db *sql.DB
	
}

func (r *PostRepository) GetPost() ([]Post, error){
	row, err := r.db.Query("SELECT id, user_id, title, content, created_at FROM posts")

	if err != nil {
		return nil, err
	}

	defer row.Close()

	var post []Post
	for row.Next() {
		var p Post

		row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
		post = append(post, p)
	}
	return post, nil
}


func (r *PostRepository) GetPostByCategory(categoryName string) ([]Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at
		FROM posts p
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ?
	`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
		posts = append(posts, p)
	}

	return posts, nil
}
func (r *PostRepository) GetPostByID(id int) (Post, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?",
		id,
	)

	var post Post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (r *PostRepository) GetPostByUser(userID int)([]Post, error) {
	rows, err := r.db.Query(`
	SELECT id, user_id, content, created_at
	FROM posts
	WHERE user_id = ? 
	`,userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan (
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post )
	}

	if err := rows.Err();err != nil {
		return nil, err
	}

	return posts, nil


}

func (r *PostRepository) CreatePost(post Post) error {
	panic("Not implemented")
}

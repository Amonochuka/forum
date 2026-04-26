package post

import "database/sql"

type PostRepository struct {
	db *sql.DB
}


type CategoryRepository struct {
	db *sql.DB
}

type UserRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *PostRepository) GetPost() ([]Post, error) {
	row, err := r.db.Query("SELECT id, user_id, title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var post []Post
	for row.Next() {
		var p Post

		err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
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

func (r *PostRepository) GetPostByUser(userID int) ([]Post, error) {
	rows, err := r.db.Query(`
	SELECT id, user_id, title, content, created_at
	FROM posts
	WHERE user_id = ? 
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostsLikedByUser(userID int) ([]Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at
		FROM posts p
		JOIN reactions re ON p.id = re.post_id
		WHERE re.user_id = ? AND re.reaction_type = 1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) CreatePost(userID int, title, content string) (int, error) {
	result, err := r.db.Exec(`
		INSERT INTO posts (user_id, title, content, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`, userID, title, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (r *PostRepository) AddPostCategory(postID, categoryID int) error {
	_, err := r.db.Exec(`
		INSERT INTO post_categories (post_id, category_id)
		VALUES (?, ?)
	`, postID, categoryID)

	return err
}

func (r *CategoryRepository) GetCategoryIDByName(name string) (int, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id FROM categories WHERE name = ?
	`, name).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CategoryRepository) GetByPostID(postID int) ([]Category, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.name
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepository) GetAllCategories() ([]Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *UserRepository) GetUsernameByID(userID int) (string, error) {
	var username string

	err := r.db.QueryRow(`
		SELECT username FROM users WHERE id = ?
	`, userID).Scan(&username)

	if err != nil {
		return "", err
	}

	return username, nil
}
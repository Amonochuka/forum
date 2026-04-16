package comment

import (
	"database/sql"
	"errors"
)

var (
    ErrNotFound   = errors.New("not found")
    ErrInvalidRef = errors.New("invalid reference")
    ErrConflict   = errors.New("conflict")
    ErrInternal   = errors.New("internal error")
)

type Repository interface {
	Create(comment *Comment) error
	GetTopLevelByPostWithReactions(postID, limit, offset int) ([]Comment, error)
	GetRepliesByParentIDWithReactions(parentID int) ([]Comment, error)
	GetByID(id int) (*Comment, error)
}

type sqliteRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepo{db: db}
}

// TODO: Return custom repository errors as opposed to database errors

func (r *sqliteRepo) GetByID(id int) (*Comment, error) {
	query := `
		SELECT
			id, 
			user_id, 
			post_id, 
			parent_id,
			content, 
			created_at,
		FROM comments
		WHERE id = ?
	`

	var c Comment
	var parentID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.UserID,
		&parentID,
		&c.PostID,
		&c.Content,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		c.ParentID = &pid
	}

	return &c, nil
}

func (r *sqliteRepo) Create(comment *Comment) error {
	query := `INSERT INTO comments (user_id, post_id, parent_id, content) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		comment.UserID,
		comment.PostID,
		comment.ParentID, // can be nil. nil becomes NULL
		comment.Content,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	comment.ID = int(id)
	return nil
}

// NOTE: Using OFFSET-based pagination doesn't scale well should you have thousands of records.
// Because it reads all records up to your offset, then discards them.
// However, for the scale of this application, that's not likely to be a problem.

func (r *sqliteRepo) GetTopLevelByPostWithReactions(postID, limit, offset int) ([]Comment, error) {
	query := `
		SELECT 
			c.id, 
			c.user_id, 
			c.post_id, 
			c.content, 
			c.created_at,
		COALESCE(SUM(CASE WHEN rx.reaction_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN rx.reaction_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c
		LEFT JOIN reactions rx ON rx.comment_id = c.id
		WHERE c.post_id = ? AND c.parent_id IS NULL
		GROUP BY c.id
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanComments(rows)
}

func (r *sqliteRepo) GetRepliesByParentIDWithReactions(parentID int) ([]Comment, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.post_id,
			c.content,
			c.created_at,
			COALESCE(SUM(CASE WHEN rx.reaction_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN rx.reaction_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c
		LEFT JOIN reactions rx ON rx.comment_id = c.id
		WHERE c.parent_id = ?
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanComments(rows)
}

func scanComments(rows *sql.Rows) ([]Comment, error) {
	var comments []Comment

	for rows.Next() {
		var c Comment

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PostID,
			&c.Content,
			&c.CreatedAt,
			&c.Likes,
			&c.Dislikes,
		)
		if err != nil {
			return nil, err
		}	

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

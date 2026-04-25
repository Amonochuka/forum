package comment

import (
	"database/sql"
	"errors"
	"strings"
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
	GetCountByPostID(postID int) (int, error)
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
		if isConstraintError(err) {
			return ErrInvalidRef
		}
		return ErrInternal
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ErrInternal
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
			u.username,
			likes,
			dislikes,
			reply_count,

			COALESCE(SUM(CASE WHEN rx.reaction_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN rx.reaction_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes

			COUNT(DISTINCT replies.id) AS reply_count

		FROM comments c

		LEFT JOIN reactions rx ON rx.comment_id = c.id
		LEFT JOIN comments replies ON replies.parent_id = c.id
		JOIN users u ON u.id = c.user_id

		WHERE c.post_id = ? AND c.parent_id IS NULL

		GROUP BY c.id

		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	comments, err := scanCommentsWithReplyCount(rows)
	if err != nil {
		return nil, ErrInternal
	}

	return comments, nil
}

func (r *sqliteRepo) GetRepliesByParentIDWithReactions(parentID int) ([]Comment, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.post_id,
			c.content,
			c.created_at,
			u.username,
			likes,
			dislikes,

			COALESCE(SUM(CASE WHEN rx.reaction_type = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN rx.reaction_type = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c
		LEFT JOIN reactions rx ON rx.comment_id = c.id
		JOIN users u ON u.id = c.user_id
		WHERE c.parent_id = ?
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	comments, err := scanComments(rows)
	if err != nil {
		return nil, ErrInternal
	}

	return comments, nil
}

func (r *sqliteRepo) GetCountByPostID(postID int) (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM comments
		WHERE post_id = ?
	`

	err := r.db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, ErrInternal
	}

	return count, nil
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
			&c.Name,
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

func scanCommentsWithReplyCount(rows *sql.Rows) ([]Comment, error) {
	var comments []Comment

	for rows.Next() {
		var c Comment

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PostID,
			&c.Content,
			&c.CreatedAt,
			&c.Name,
			&c.Likes,
			&c.Dislikes,
			&c.ReplyCount,
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

func isConstraintError(err error) bool {
	var sqliteErr interface{ Error() string }
	if errors.As(err, &sqliteErr) {
		return strings.Contains(sqliteErr.Error(), "FOREIGN KEY constraint failed")
	}
	return false
}

package reaction

import (
	"database/sql"
	"errors"
	"fmt"
)

type ReactionRepository struct {
	DB *sql.DB
}

func (r *ReactionRepository) AddReaction(reaction *Reaction) error {
	_, err := r.DB.Exec(`
	INSERT INTO reactions(user_id, post_id, comment_id, reaction_type)
	VALUES(?,?,?,?)`,
		reaction.UserID,
		reaction.PostID,
		reaction.CommentID,
		reaction.Type,
	)
	if err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}
	return nil
}

func (r *ReactionRepository) UpdateReaction(reaction *Reaction) error {
	if reaction.PostID != nil {
		_, err := r.DB.Exec(`
		UPDATE reactions SET reaction_type = ? WHERE user_id = ? AND post_id = ?`,
			reaction.Type,
			reaction.UserID,
			*reaction.PostID,
		)

		if err != nil {
			return fmt.Errorf("update reaction (post): %w", err)
		}

		return nil
	}

	_, err := r.DB.Exec(`
	UPDATE reactions SET reaction_type = ? WHERE user_id = ? AND comment_id = ?`,
		reaction.Type,
		reaction.UserID,
		*reaction.CommentID,
	)

	if err != nil {
		return fmt.Errorf("update reaction (comment): %w", err)
	}

	return nil
}

func (r *ReactionRepository) DeleteReaction(userID int, postID *int, commentID *int) error {
	if postID != nil {
		res, err := r.DB.Exec(`
		DELETE FROM reactions WHERE user_id = ? AND post_id = ?`,
			userID,
			*postID,
		)

		if err != nil {
			return fmt.Errorf("delete reaction (post): %w", err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("delete reaction rows affected (post): %w", err)
		}

		if rows == 0 {
			return errors.New("reaction not found")
		}

		return nil
	}

	res, err := r.DB.Exec(`
	DELETE FROM reactions WHERE user_id = ? AND comment_id = ?`,
		userID,
		*commentID,
	)

	if err != nil {
		return fmt.Errorf("delete reaction (comment): %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete reaction rows affected (comment): %w", err)
	}

	if rows == 0 {
		return errors.New("reaction not found")
	}

	return nil
}

func (r *ReactionRepository) GetPostReactions(postID int) ([]*Reaction, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, post_id, comment_id, reaction_type, created_at
	FROM reactions
	WHERE post_id = ?`, postID)

	if err != nil {
		return nil, fmt.Errorf("get post reactions: %w", err)
	}
	defer rows.Close()

	var reactions []*Reaction
	for rows.Next() {
		reaction := &Reaction{}
		err := rows.Scan(
			&reaction.ID,
			&reaction.UserID,
			&reaction.PostID,
			&reaction.CommentID,
			&reaction.Type,
			&reaction.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan post reactions: %w", err)
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (r *ReactionRepository) GetCommentReactions(commentID int) ([]*Reaction, error) {
	rows, err := r.DB.Query(`
	SELECT id, user_id, post_id, comment_id, reaction_type, created_at
	FROM reactions
	WHERE comment_id = ?`, commentID)

	if err != nil {
		return nil, fmt.Errorf("get comment reactions: %w", err)
	}
	defer rows.Close()

	var reactions []*Reaction
	for rows.Next() {
		reaction := &Reaction{}

		err := rows.Scan(
			&reaction.ID,
			&reaction.UserID,
			&reaction.PostID,
			&reaction.CommentID,
			&reaction.Type,
			&reaction.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan comment reactions: %w", err)
		}

		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func (r *ReactionRepository) GetUserReaction(userID int, postID *int, commentID *int) (*Reaction, error) {
	if postID != nil {
		row := r.DB.QueryRow(`
		SELECT id, user_id, post_id, comment_id, reaction_type, created_at
		FROM reactions
		WHERE user_id = ? AND post_id = ?`, userID, *postID)
		reaction := &Reaction{}
		err := row.Scan(
			&reaction.ID,
			&reaction.UserID,
			&reaction.PostID,
			&reaction.CommentID,
			&reaction.Type,
			&reaction.CreatedAt,
		)

		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			return nil, fmt.Errorf("get user reaction (post): %w", err)
		}

		return reaction, nil
	}

	row := r.DB.QueryRow(`
	SELECT id, user_id, post_id, comment_id, reaction_type, created_at
	FROM reactions
	WHERE user_id = ? AND comment_id = ?`, userID, *commentID)

	reaction := &Reaction{}
	err := row.Scan(
		&reaction.ID,
		&reaction.UserID,
		&reaction.PostID,
		&reaction.CommentID,
		&reaction.Type,
		&reaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user reaction (comment): %w", err)
	}
	return reaction, nil
}

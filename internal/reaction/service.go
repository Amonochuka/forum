package reaction

import "errors"

type ReactionService struct {
	Repo *ReactionRepository
}

func (s *ReactionService) React(r *Reaction) error {
	//validate reaction
	if (r.PostID == nil && r.CommentID == nil) ||
		(r.PostID != nil && r.CommentID != nil) {
		return errors.New("reaction must belong to either post or comment")
	}

	//check existing reaction
	existing, err := s.Repo.GetUserReaction(r.UserID, r.PostID, r.CommentID)
	if err != nil {
		return err
	}

	//insert if there is no reaction
	if existing == nil {
		return s.Repo.AddReaction(r)
	}

	//delete if reaction is same as before
	if existing.Type == r.Type {
		return s.Repo.DeleteReaction(r.UserID, r.PostID, r.CommentID)
	}

	//update of the reaction is different
	return s.Repo.UpdateReaction(r)
}

func (s *ReactionService) GetPostReactions(postID int) ([]*Reaction, error) {
	return s.Repo.GetPostReactions(postID)
}

func (s *ReactionService) GetCommentReactions(commentID int) ([]*Reaction, error) {
	return s.Repo.GetCommentReactions(commentID)
}
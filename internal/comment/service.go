package comment

import "errors"

var (
	ErrEmptyContent          = errors.New("content cannot be empty")
	ErrContentTooLong        = errors.New("content too long")
	ErrInvalidParentID       = errors.New("invalid parent id")
	ErrNestedReplyNotAllowed = errors.New("nested replies not allowed")
	ErrInternalServerError   = errors.New("internal server error")
)

type Service interface {
	CreateComment(userID, postID int, content string, parentID *int) error
	GetTopLevelComments(postID, page int) ([]Comment, error)
	GetReplies(parentID int) ([]Comment, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateComment(userID, postID int, content string, parentID *int) error {
	if content == "" {
		return ErrEmptyContent
	}

	if len(content) > 1000 {
		return ErrContentTooLong
	}

	// Prevent repling to a reply
	if parentID != nil {
		parent, err := s.repo.GetByID(*parentID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrInvalidParentID
			}
			return ErrInternalServerError
		}

		if parent.ParentID != nil {
			return ErrNestedReplyNotAllowed
		}
	}

	comment := &Comment{
		UserID:   userID,
		PostID:   postID,
		ParentID: parentID,
		Content:  content,
	}

	if err := s.repo.Create(comment); err != nil {
		if errors.Is(err, ErrInvalidRef) {
			return ErrInvalidParentID
		}
		return ErrInternalServerError
	}

	return nil
}

func (s *service) GetTopLevelComments(postID, page int) ([]Comment, error) {
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	comments, err := s.repo.GetTopLevelByPostWithReactions(postID, limit, offset)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return comments, nil
}

func (s *service) GetReplies(parentID int) ([]Comment, error) {
	if parentID <= 0 {
		return nil, ErrInvalidParentID
	}

	replies, err := s.repo.GetRepliesByParentIDWithReactions(parentID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return replies, nil
}

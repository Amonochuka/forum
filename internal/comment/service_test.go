package comment

import (
	"errors"
	"testing"
)

// MockRepo implements the comment.Repository interface
type MockRepo struct {
	CreateFunc                        func(comment *Comment) error
	GetTopLevelByPostWithReactionsFunc func(postID, limit, offset int) ([]Comment, error)
	GetRepliesByParentIDWithReactionsFunc func(parentID int) ([]Comment, error)
	GetByIDFunc                       func(id int) (*Comment, error)
	GetCountByPostIDFunc              func(postID int) (int, error)
}

func (m *MockRepo) Create(c *Comment) error { return m.CreateFunc(c) }
func (m *MockRepo) GetTopLevelByPostWithReactions(p, l, o int) ([]Comment, error) {
	return m.GetTopLevelByPostWithReactionsFunc(p, l, o)
}
func (m *MockRepo) GetRepliesByParentIDWithReactions(id int) ([]Comment, error) {
	return m.GetRepliesByParentIDWithReactionsFunc(id)
}
func (m *MockRepo) GetByID(id int) (*Comment, error) { return m.GetByIDFunc(id) }
func (m *MockRepo) GetCountByPostID(id int) (int, error) { return m.GetCountByPostIDFunc(id) }

// --- Tests ---

func TestCreateComment_PreventNestedReplies(t *testing.T) {
	// Setup: Parent is a reply (it has a ParentID), so we shouldn't be allowed to reply to it
	parentID := 10
	mock := &MockRepo{
		GetByIDFunc: func(id int) (*Comment, error) {
			// Simulating that the parent itself is already a reply
			return &Comment{ID: id, ParentID: new(int)}, nil
		},
	}

	service := NewService(mock)

	// Execute: Try to create a comment with this parent
	_, err := service.CreateComment(1, 101, "Trying to nest", &parentID)

	// Verify: Should fail with ErrNestedReplyNotAllowed
	if !errors.Is(err, ErrNestedReplyNotAllowed) {
		t.Errorf("Expected ErrNestedReplyNotAllowed, got %v", err)
	}
}

func TestCreateComment_Success(t *testing.T) {
	// Setup: Repository should accept the new comment
	mock := &MockRepo{
		CreateFunc: func(c *Comment) error { return nil },
	}

	service := NewService(mock)

	// Execute
	_, err := service.CreateComment(1, 101, "Valid comment", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
}
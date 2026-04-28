package reaction

import (
	"errors"
	"testing"
)

/*
MockReactionRepo implements FULL Repository interface
so it compiles with Go interface rules
*/
type MockReactionRepo struct {
	GetUserReactionFunc            func(userID int, postID *int, commentID *int) (*Reaction, error)
	AddReactionFunc                func(r *Reaction) error
	UpdateReactionFunc             func(r *Reaction) error
	DeleteReactionFunc             func(userID int, postID *int, commentID *int) error

	GetPostReactionsFunc           func(postID int) ([]*Reaction, error)
	GetCommentReactionsFunc        func(commentID int) ([]*Reaction, error)
	GetPostReactionCountsFunc      func(postID int) (int, int, error)
	GetCommentReactionCountsFunc   func(commentID int) (int, int, error)
}

// ---------------- CORE METHODS ----------------

func (m *MockReactionRepo) GetUserReaction(userID int, postID *int, commentID *int) (*Reaction, error) {
	if m.GetUserReactionFunc == nil {
		return nil, nil
	}
	return m.GetUserReactionFunc(userID, postID, commentID)
}

func (m *MockReactionRepo) AddReaction(r *Reaction) error {
	if m.AddReactionFunc == nil {
		return nil
	}
	return m.AddReactionFunc(r)
}

func (m *MockReactionRepo) UpdateReaction(r *Reaction) error {
	if m.UpdateReactionFunc == nil {
		return nil
	}
	return m.UpdateReactionFunc(r)
}

func (m *MockReactionRepo) DeleteReaction(userID int, postID *int, commentID *int) error {
	if m.DeleteReactionFunc == nil {
		return nil
	}
	return m.DeleteReactionFunc(userID, postID, commentID)
}

// ---------------- OPTIONAL METHODS ----------------

func (m *MockReactionRepo) GetPostReactions(postID int) ([]*Reaction, error) {
	if m.GetPostReactionsFunc == nil {
		return nil, nil
	}
	return m.GetPostReactionsFunc(postID)
}

func (m *MockReactionRepo) GetCommentReactions(commentID int) ([]*Reaction, error) {
	if m.GetCommentReactionsFunc == nil {
		return nil, nil
	}
	return m.GetCommentReactionsFunc(commentID)
}

func (m *MockReactionRepo) GetPostReactionCounts(postID int) (int, int, error) {
	if m.GetPostReactionCountsFunc == nil {
		return 0, 0, nil
	}
	return m.GetPostReactionCountsFunc(postID)
}

func (m *MockReactionRepo) GetCommentReactionCounts(commentID int) (int, int, error) {
	if m.GetCommentReactionCountsFunc == nil {
		return 0, 0, nil
	}
	return m.GetCommentReactionCountsFunc(commentID)
}

// -------------------- TESTS --------------------

func TestReact_InvalidTarget(t *testing.T) {
	service := &ReactionService{
		Repo: &MockReactionRepo{},
	}

	err := service.React(&Reaction{
		UserID: 1,
		Type:   1,
	})

	if err == nil {
		t.Fatalf("expected error for invalid reaction target")
	}
}

func TestReact_AddsReaction_WhenNoneExists(t *testing.T) {
	postID := 1
	addCalled := false

	mock := &MockReactionRepo{
		GetUserReactionFunc: func(int, *int, *int) (*Reaction, error) {
			return nil, nil
		},
		AddReactionFunc: func(r *Reaction) error {
			addCalled = true
			return nil
		},
	}

	service := &ReactionService{Repo: mock}

	err := service.React(&Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !addCalled {
		t.Fatalf("expected AddReaction to be called")
	}
}

func TestReact_TogglesReaction_WhenSameType(t *testing.T) {
	postID := 1
	deleteCalled := false

	mock := &MockReactionRepo{
		GetUserReactionFunc: func(int, *int, *int) (*Reaction, error) {
			return &Reaction{Type: 1}, nil
		},
		DeleteReactionFunc: func(int, *int, *int) error {
			deleteCalled = true
			return nil
		},
	}

	service := &ReactionService{Repo: mock}

	err := service.React(&Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !deleteCalled {
		t.Fatalf("expected DeleteReaction to be called")
	}
}

func TestReact_UpdatesReaction_WhenDifferentType(t *testing.T) {
	postID := 1
	updateCalled := false

	mock := &MockReactionRepo{
		GetUserReactionFunc: func(int, *int, *int) (*Reaction, error) {
			return &Reaction{Type: 1}, nil
		},
		UpdateReactionFunc: func(r *Reaction) error {
			updateCalled = true
			if r.Type != -1 {
				t.Fatalf("expected reaction type -1, got %d", r.Type)
			}
			return nil
		},
	}

	service := &ReactionService{Repo: mock}

	err := service.React(&Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   -1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !updateCalled {
		t.Fatalf("expected UpdateReaction to be called")
	}
}

func TestReact_PropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo error")

	mock := &MockReactionRepo{
		GetUserReactionFunc: func(int, *int, *int) (*Reaction, error) {
			return nil, expectedErr
		},
	}

	service := &ReactionService{Repo: mock}

	postID := 1
	err := service.React(&Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
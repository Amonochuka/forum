package comment

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"forum/internal/shared/middleware"
	"forum/internal/user"
)

type stubCommentService struct {
	createCommentFunc func(userID, postID int, content string, parentID *int) (*Comment, error)
}

func (s *stubCommentService) CreateComment(userID, postID int, content string, parentID *int) (*Comment, error) {
	return s.createCommentFunc(userID, postID, content, parentID)
}

func (s *stubCommentService) GetTopLevelComments(postID, page int) ([]Comment, int, error) {
	return nil, 0, nil
}

func (s *stubCommentService) GetReplies(parentID int) ([]Comment, error) {
	return nil, nil
}

type stubUserService struct {
	getByIDFunc func(id int) (*user.Profile, error)
}

func (s *stubUserService) GetByID(id int) (*user.Profile, error) {
	return s.getByIDFunc(id)
}

func TestHandler_CreateComment_Success(t *testing.T) {
	handler := NewHandler(
		&stubCommentService{
			createCommentFunc: func(userID, postID int, content string, parentID *int) (*Comment, error) {
				return &Comment{
					ID:        7,
					UserID:    userID,
					PostID:    postID,
					Content:   content,
					CreatedAt: time.Now(),
				}, nil
			},
		},
		&stubUserService{
			getByIDFunc: func(id int) (*user.Profile, error) {
				return &user.Profile{ID: id, Username: "francis"}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/posts/12/comments", strings.NewReader("content=Hello+there"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = middleware.WithUserID(req, 1)
	req.SetPathValue("id", "12")

	rr := httptest.NewRecorder()
	handler.CreateComment(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp struct {
		Comment CommentView `json:"comment"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Comment.AuthorName != "francis" {
		t.Fatalf("expected author name francis, got %q", resp.Comment.AuthorName)
	}
	if resp.Comment.Body != "Hello there" {
		t.Fatalf("expected body %q, got %q", "Hello there", resp.Comment.Body)
	}
}

func TestHandler_CreateReply_EmptyContent(t *testing.T) {
	handler := NewHandler(
		&stubCommentService{
			createCommentFunc: func(userID, postID int, content string, parentID *int) (*Comment, error) {
				return nil, ErrEmptyContent
			},
		},
		&stubUserService{
			getByIDFunc: func(id int) (*user.Profile, error) {
				return &user.Profile{ID: id, Username: "francis"}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/comments/5/replies", strings.NewReader("content=+++&post_id=12"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = middleware.WithUserID(req, 1)
	req.SetPathValue("id", "5")

	rr := httptest.NewRecorder()
	handler.CreateReply(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "reply cannot be empty") {
		t.Fatalf("expected empty reply message, got %q", rr.Body.String())
	}
}

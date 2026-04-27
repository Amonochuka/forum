package post

import (
	"errors"
	"testing"
)

// --- Mocks ---

type MockPostRepo struct {
	GetPostFunc            func() ([]Post, error)
	GetPostByIDFunc        func(id int) (Post, error)
	CreatePostFunc         func(userID int, title, content string) (int, error)
	GetPostByCategoryFunc  func(cat string) ([]Post, error)
	AddPostCategoryFunc    func(postID, catID int) error
}

func (m *MockPostRepo) GetPost() ([]Post, error)                         { return m.GetPostFunc() }
func (m *MockPostRepo) GetPostByID(id int) (Post, error)                { return m.GetPostByIDFunc(id) }
func (m *MockPostRepo) GetPostByCategory(c string) ([]Post, error)      { return m.GetPostByCategoryFunc(c) }
func (m *MockPostRepo) CreatePost(u int, t, c string) (int, error)      { return m.CreatePostFunc(u, t, c) }
func (m *MockPostRepo) AddPostCategory(p, c int) error                  { return m.AddPostCategoryFunc(p, c) }
func (m *MockPostRepo) GetPostByUser(id int) ([]Post, error)            { return nil, nil }
func (m *MockPostRepo) GetPostsLikedByUser(id int) ([]Post, error)      { return nil, nil }

type MockUserRepo struct {
	GetUsernameByIDFunc func(id int) (string, error)
}
func (m *MockUserRepo) GetUsernameByID(id int) (string, error) { return m.GetUsernameByIDFunc(id) }

type MockCategoryRepo struct {
	GetByPostIDFunc         func(id int) ([]Category, error)
	GetCategoryIDByNameFunc func(name string) (int, error)
}
func (m *MockCategoryRepo) GetByPostID(id int) ([]Category, error)     { return m.GetByPostIDFunc(id) }
func (m *MockCategoryRepo) GetCategoryIDByName(n string) (int, error) { return m.GetCategoryIDByNameFunc(n) }
func (m *MockCategoryRepo) GetAllCategories() ([]Category, error)      { return nil, nil }

type MockReactionRepo struct {
	GetPostReactionCountsFunc func(id int) (int, int, error)
}
func (m *MockReactionRepo) GetPostReactionCounts(id int) (int, int, error) { 
	return m.GetPostReactionCountsFunc(id) 
}

// --- Tests ---

func TestGetPosts_AggregationLogic(t *testing.T) {
	// Setup: We want to test if buildPostResponse correctly joins user and category data
	mockPost := &MockPostRepo{
		GetPostFunc: func() ([]Post, error) {
			return []Post{{ID: 1, Title: "Test", UserID: 5}}, nil
		},
	}
	mockUser := &MockUserRepo{
		GetUsernameByIDFunc: func(id int) (string, error) { return "Gopher", nil },
	}
	mockCat := &MockCategoryRepo{
		GetByPostIDFunc: func(id int) ([]Category, error) {
			return []Category{{ID: 1, Name: "Golang"}}, nil
		},
	}
	mockReact := &MockReactionRepo{
		GetPostReactionCountsFunc: func(id int) (int, int, error) {
			return 10, 2, nil
		},
	}

	service := NewPostService(mockPost, mockCat, mockUser, mockReact)

	// Execute
	results, err := service.GetPosts("", "", "")

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 post, got %d", len(results))
	}
	if results[0].Username != "Gopher" {
		t.Errorf("Expected username Gopher, got %s", results[0].Username)
	}
	if results[0].Likes != 10 {
		t.Errorf("Expected 10 likes, got %d", results[0].Likes)
	}
}

func TestCreatePost_Validation(t *testing.T) {
	service := NewPostService(nil, nil, nil, nil)

	// Execute: Empty title/content
	err := service.CreatePost(1, "", "", nil)

	// Verify
	if !errors.Is(err, ErrEmptyContent) {
		t.Errorf("Expected ErrEmptyContent, got %v", err)
	}
}
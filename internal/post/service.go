package post

import (
	"errors"
	"strconv"
)

type PostService struct {
	postRepo     *PostRepository
	userRepo     *UserRepository
	categoryRepo *CategoryRepository
}

var ErrEmptyContent = errors.New("content cannot be empty")

func NewPostService(repo *PostRepository, userRepo *UserRepository, categoryRepo *CategoryRepository) *PostService {
	return &PostService{postRepo: repo,userRepo:userRepo,categoryRepo: categoryRepo}
}

func (s *PostService) buildPostResponse(post Post) (PostResponse, error) {
	username, err := s.userRepo.GetUsernameByID(post.UserID)
	if err != nil {
		return PostResponse{}, err
	}

	categories, err := s.categoryRepo.GetByPostID(post.ID)
	if err != nil {
		return PostResponse{}, err
	}

	return PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Username:  username,
		Category:  categories,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (s *PostService) GetPosts(category, user string) ([]PostResponse, error) {
	var posts []Post
	var err error

	// Apply filters
	if category != "" {
		posts, err = s.postRepo.GetPostByCategory(category)
	} else if user != "" {
		userID, err := strconv.Atoi(user)
		if err != nil {
			return nil, err
		}
		posts, err = s.postRepo.GetPostByUser(userID)
	} else {
		posts, err = s.postRepo.GetPost()
	}

	if err != nil {
		return nil, err
	}

	// Build response
	var result []PostResponse

	for _, post := range posts {
		resp, err := s.buildPostResponse(post)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}

	return result, nil
}

func (s *PostService) GetPostByID(id int) (PostResponse, error) {
	post, err := s.postRepo.GetPostByID(id)
	if err != nil {
		return PostResponse{}, err
	}

	return s.buildPostResponse(post)
}

func (s *PostService) CreatePost(userID int, title, content string, categories []string) error {
	// Basic validation
	if title == "" || content == "" {
		return ErrEmptyContent
	}

	// 1. Create post
	postID, err := s.postRepo.CreatePost(userID, title, content)
	if err != nil {
		return err
	}

	// 2. Attach categories
	for _, catName := range categories {
		// get category ID
		catID, err := s.categoryRepo.GetCategoryIDByName(catName)
		if err != nil {
			return err
		}

		// link post ↔ category
		err = s.postRepo.AddPostCategory(postID, catID)
		if err != nil {
			return err
		}
	}

	return nil
}

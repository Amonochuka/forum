package post

import "strconv"

type PostService struct {
	postRepo *PostRepository
}

func (s *PostService) buildPostResponse(post Post) (PostResponse, error) {
	username, err := s.userRepo.GetUsernameByID(post.UserID)
	if err != nil {
		return PostResponse{}, err
	}

	categories, err := s.categoryRepo.GetByPostId(post.ID)
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
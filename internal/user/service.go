package user

type Service interface {
	GetByID(id int) (*Profile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(id int) (*Profile, error) {
	return s.repo.GetByID(id)
}
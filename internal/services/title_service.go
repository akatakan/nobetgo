package services

import (
	"github.com/akatakan/nobetgo/internal/core"
)

type TitleRepositoryInterface interface {
	Create(title *core.Title) error
	GetByID(id uint) (*core.Title, error)
	GetByName(name string) (*core.Title, error)
	List() ([]core.Title, error)
	Update(title *core.Title) error
	Delete(id uint) error
}

type TitleService struct {
	repo TitleRepositoryInterface
}

func NewTitleService(repo TitleRepositoryInterface) *TitleService {
	return &TitleService{repo: repo}
}

func (s *TitleService) CreateTitle(title *core.Title) error {
	return s.repo.Create(title)
}

func (s *TitleService) GetTitleByID(id uint) (*core.Title, error) {
	return s.repo.GetByID(id)
}

func (s *TitleService) GetAllTitles() ([]core.Title, error) {
	return s.repo.List()
}

func (s *TitleService) UpdateTitle(title *core.Title) error {
	return s.repo.Update(title)
}

func (s *TitleService) DeleteTitle(id uint) error {
	return s.repo.Delete(id)
}

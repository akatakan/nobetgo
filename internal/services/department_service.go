package services

import (
	"github.com/akatakan/nobetgo/internal/core"
)

type DepartmentRepositoryInterface interface {
	Create(department *core.Department) error
	GetByID(id uint) (*core.Department, error)
	GetByName(name string) (*core.Department, error)
	List() ([]core.Department, error)
	Update(department *core.Department) error
	Delete(id uint) error
}

type DepartmentService struct {
	repo DepartmentRepositoryInterface
}

func NewDepartmentService(repo DepartmentRepositoryInterface) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) CreateDepartment(department *core.Department) error {
	return s.repo.Create(department)
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*core.Department, error) {
	return s.repo.GetByID(id)
}

func (s *DepartmentService) GetAllDepartments() ([]core.Department, error) {
	return s.repo.List()
}

func (s *DepartmentService) UpdateDepartment(department *core.Department) error {
	return s.repo.Update(department)
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	return s.repo.Delete(id)
}

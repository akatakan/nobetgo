package services

import (
	"github.com/akatakan/nobetgo/internal/core"
)

type EmployeeRepositoryInterface interface {
	Create(employee *core.Employee) error
	GetByID(id uint) (*core.Employee, error)
	List() ([]core.Employee, error)
	Update(employee *core.Employee) error
	Delete(id uint) error
}

type EmployeeService struct {
	repo EmployeeRepositoryInterface
}

func NewEmployeeService(repo EmployeeRepositoryInterface) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) CreateEmployee(employee *core.Employee) error {
	return s.repo.Create(employee)
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*core.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) GetAllEmployees() ([]core.Employee, error) {
	return s.repo.List()
}

func (s *EmployeeService) UpdateEmployee(employee *core.Employee) error {
	return s.repo.Update(employee)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	return s.repo.Delete(id)
}

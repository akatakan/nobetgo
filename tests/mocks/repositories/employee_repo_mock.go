package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEmployeeRepository is a mock implementation of EmployeeRepositoryInterface
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(employee *core.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(id uint) (*core.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) List() ([]core.Employee, error) {
	args := m.Called()
	return args.Get(0).([]core.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(employee *core.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) ListByDepartment(departmentID uint) ([]core.Employee, error) {
	args := m.Called(departmentID)
	return args.Get(0).([]core.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

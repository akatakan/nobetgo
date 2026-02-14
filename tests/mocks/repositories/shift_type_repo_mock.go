package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/stretchr/testify/mock"
)

// MockShiftTypeRepository is a mock implementation of ShiftTypeRepositoryInterface
type MockShiftTypeRepository struct {
	mock.Mock
}

func (m *MockShiftTypeRepository) Create(shiftType *core.ShiftType) error {
	args := m.Called(shiftType)
	return args.Error(0)
}

func (m *MockShiftTypeRepository) GetByID(id uint) (*core.ShiftType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.ShiftType), args.Error(1)
}

func (m *MockShiftTypeRepository) List() ([]core.ShiftType, error) {
	args := m.Called()
	return args.Get(0).([]core.ShiftType), args.Error(1)
}

func (m *MockShiftTypeRepository) Update(shiftType *core.ShiftType) error {
	args := m.Called(shiftType)
	return args.Error(0)
}

func (m *MockShiftTypeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

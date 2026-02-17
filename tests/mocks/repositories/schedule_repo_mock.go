package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRepository is a mock implementation of ScheduleRepositoryInterface
type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) Create(schedule *core.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) Update(schedule *core.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetCombinedSchedule(month int, year int) ([]core.Schedule, error) {
	args := m.Called(month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]core.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) DeleteByMonthYear(month int, year int) error {
	args := m.Called(month, year)
	return args.Error(0)
}

func (m *MockScheduleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetByID(id uint) (*core.Schedule, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Schedule), args.Error(1)
}

package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/stretchr/testify/mock"
)

// MockTimeEntryRepository is a mock implementation of TimeEntryRepositoryInterface.
type MockTimeEntryRepository struct {
	mock.Mock
}

func (m *MockTimeEntryRepository) Create(entry *core.TimeEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockTimeEntryRepository) Update(entry *core.TimeEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockTimeEntryRepository) GetByID(id uint) (*core.TimeEntry, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTimeEntryRepository) GetOpenEntry(employeeID uint) (*core.TimeEntry, error) {
	args := m.Called(employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepository) ListByEmployee(employeeID uint, start, end time.Time) ([]core.TimeEntry, error) {
	args := m.Called(employeeID, start, end)
	return args.Get(0).([]core.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepository) ListByDepartment(departmentID uint, start, end time.Time) ([]core.TimeEntry, error) {
	args := m.Called(departmentID, start, end)
	return args.Get(0).([]core.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepository) ListByDateRange(start, end time.Time) ([]core.TimeEntry, error) {
	args := m.Called(start, end)
	return args.Get(0).([]core.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepository) ListByStatus(status string, start, end time.Time) ([]core.TimeEntry, error) {
	args := m.Called(status, start, end)
	return args.Get(0).([]core.TimeEntry), args.Error(1)
}

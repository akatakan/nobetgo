package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/stretchr/testify/mock"
)

type MockLeaveRepository struct {
	mock.Mock
}

func (m *MockLeaveRepository) Create(leave *core.Leave) error {
	args := m.Called(leave)
	return args.Error(0)
}

func (m *MockLeaveRepository) Update(leave *core.Leave) error {
	args := m.Called(leave)
	return args.Error(0)
}

func (m *MockLeaveRepository) GetByID(id uint) (*core.Leave, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Leave), args.Error(1)
}

func (m *MockLeaveRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLeaveRepository) ListByEmployee(employeeID uint, start, end time.Time) ([]core.Leave, error) {
	args := m.Called(employeeID, start, end)
	return args.Get(0).([]core.Leave), args.Error(1)
}

func (m *MockLeaveRepository) ListByDepartment(departmentID uint, start, end time.Time) ([]core.Leave, error) {
	args := m.Called(departmentID, start, end)
	return args.Get(0).([]core.Leave), args.Error(1)
}

func (m *MockLeaveRepository) ListByStatus(status string) ([]core.Leave, error) {
	args := m.Called(status)
	return args.Get(0).([]core.Leave), args.Error(1)
}

func (m *MockLeaveRepository) HasOverlap(employeeID uint, start, end time.Time, excludeID uint) (bool, error) {
	args := m.Called(employeeID, start, end, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLeaveRepository) CreateLeaveType(lt *core.LeaveType) error {
	args := m.Called(lt)
	return args.Error(0)
}

func (m *MockLeaveRepository) UpdateLeaveType(lt *core.LeaveType) error {
	args := m.Called(lt)
	return args.Error(0)
}

func (m *MockLeaveRepository) GetLeaveTypeByID(id uint) (*core.LeaveType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.LeaveType), args.Error(1)
}

func (m *MockLeaveRepository) ListLeaveTypes() ([]core.LeaveType, error) {
	args := m.Called()
	return args.Get(0).([]core.LeaveType), args.Error(1)
}

func (m *MockLeaveRepository) DeleteLeaveType(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLeaveRepository) GetBalance(employeeID uint, leaveTypeID uint, year int) (*core.LeaveBalance, error) {
	args := m.Called(employeeID, leaveTypeID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.LeaveBalance), args.Error(1)
}

func (m *MockLeaveRepository) GetAllBalances(employeeID uint, year int) ([]core.LeaveBalance, error) {
	args := m.Called(employeeID, year)
	return args.Get(0).([]core.LeaveBalance), args.Error(1)
}

func (m *MockLeaveRepository) UpsertBalance(balance *core.LeaveBalance) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *MockLeaveRepository) ListPaginated(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) ([]core.Leave, int64, error) {
	args := m.Called(params, employeeID, departmentID, start, end)
	return args.Get(0).([]core.Leave), args.Get(1).(int64), args.Error(2)
}

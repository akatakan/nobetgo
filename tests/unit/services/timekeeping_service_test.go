package services_test

import (
	"testing"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	mocks "github.com/akatakan/nobetgo/tests/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock Repo
type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) Create(attendance *core.Attendance) error {
	args := m.Called(attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) Update(attendance *core.Attendance) error {
	args := m.Called(attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetByID(id uint) (*core.Attendance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Attendance), args.Error(1)
}

func (m *MockAttendanceRepo) GetCombinedReport(month int, year int) ([]core.Attendance, error) {
	args := m.Called(month, year)
	return args.Get(0).([]core.Attendance), args.Error(1)
}

func TestLogAttendance_Overtime(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockScheduleRepo := new(mocks.MockScheduleRepository) // Use existing mock package
	service := services.NewTimekeepingService(mockRepo, mockScheduleRepo)

	// Mock GetByID for overtime calculation
	schedule := &core.Schedule{
		Model: gorm.Model{ID: 1},
		ShiftType: core.ShiftType{
			StartTime: "08:00",
			EndTime:   "16:00", // 8 hours duration
		},
	}
	mockScheduleRepo.On("GetByID", uint(1)).Return(schedule, nil)

	// Case 1: 10 hours work (2 hours overtime)
	startTime := time.Date(2025, 2, 1, 8, 0, 0, 0, time.UTC)
	endTime := startTime.Add(10 * time.Hour)

	req := core.AttendanceRequest{
		ScheduleID:      1,
		ActualStartTime: startTime,
		ActualEndTime:   endTime,
		Notes:           "Overtime",
	}

	mockRepo.On("Create", mock.MatchedBy(func(a *core.Attendance) bool {
		return a.IsOvertime == true && a.OvertimeHours == 2.0
	})).Return(nil)

	attendance, err := service.LogAttendance(req)

	assert.NoError(t, err)
	assert.NotNil(t, attendance)
	assert.True(t, attendance.IsOvertime)
	assert.Equal(t, 2.0, attendance.OvertimeHours)

	mockRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestLogAttendance_Normal(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockRepo, mockScheduleRepo)

	// Mock GetByID
	schedule := &core.Schedule{
		Model: gorm.Model{ID: 1},
		ShiftType: core.ShiftType{
			StartTime: "08:00",
			EndTime:   "16:00", // 8 hours
		},
	}
	mockScheduleRepo.On("GetByID", uint(1)).Return(schedule, nil)

	// Case 2: 8 hours work (No overtime)
	startTime := time.Date(2025, 2, 1, 8, 0, 0, 0, time.UTC)
	endTime := startTime.Add(8 * time.Hour)

	req := core.AttendanceRequest{
		ScheduleID:      1,
		ActualStartTime: startTime,
		ActualEndTime:   endTime,
		Notes:           "Normal shift",
	}

	mockRepo.On("Create", mock.MatchedBy(func(a *core.Attendance) bool {
		return a.IsOvertime == false && a.OvertimeHours == 0.0
	})).Return(nil)

	attendance, err := service.LogAttendance(req)

	assert.NoError(t, err)
	assert.False(t, attendance.IsOvertime)

	mockRepo.AssertExpectations(t)
}

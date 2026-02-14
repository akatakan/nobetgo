package services_test

import (
	"testing"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repo
type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) Create(attendance *core.Attendance) error {
	args := m.Called(attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetCombinedReport(month int, year int) ([]core.Attendance, error) {
	args := m.Called(month, year)
	return args.Get(0).([]core.Attendance), args.Error(1)
}

func TestLogAttendance_Overtime(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	service := services.NewTimekeepingService(mockRepo)

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
}

func TestLogAttendance_Normal(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	service := services.NewTimekeepingService(mockRepo)

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

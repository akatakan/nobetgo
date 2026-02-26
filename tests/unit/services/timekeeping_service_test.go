package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	mocks "github.com/akatakan/nobetgo/tests/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClockIn_Success(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	// No open entry exists
	mockTimeEntryRepo.On("GetOpenEntry", uint(1)).Return(nil, errors.New("not found"))
	mockTimeEntryRepo.On("Create", mock.AnythingOfType("*core.TimeEntry")).Return(nil)

	req := core.ClockInRequest{
		EmployeeID: 1,
		Notes:      "Sabah girişi",
	}

	entry, err := service.ClockIn(req)

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, uint(1), entry.EmployeeID)
	assert.Equal(t, "auto", entry.Source)
	assert.Equal(t, "pending", entry.Status)
	assert.Nil(t, entry.ClockOut)

	mockTimeEntryRepo.AssertExpectations(t)
}

func TestClockIn_AlreadyOpen(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	// Open entry exists
	existing := &core.TimeEntry{EmployeeID: 1, ClockIn: time.Now()}
	existing.ID = 5
	mockTimeEntryRepo.On("GetOpenEntry", uint(1)).Return(existing, nil)

	req := core.ClockInRequest{EmployeeID: 1}
	entry, err := service.ClockIn(req)

	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Contains(t, err.Error(), "açık bir giriş kaydı var")
}

func TestClockOut_Success(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	// Open entry exists — clocked in 8 hours ago
	clockIn := time.Now().Add(-8 * time.Hour)
	existing := &core.TimeEntry{
		EmployeeID: 1,
		ClockIn:    clockIn,
		EntryType:  "normal",
		Source:     "auto",
		Status:     "pending",
	}
	existing.ID = 10

	mockTimeEntryRepo.On("GetOpenEntry", uint(1)).Return(existing, nil)
	mockTimeEntryRepo.On("Update", mock.AnythingOfType("*core.TimeEntry")).Return(nil)

	req := core.ClockOutRequest{EmployeeID: 1}
	entry, err := service.ClockOut(req)

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.NotNil(t, entry.ClockOut)

	mockTimeEntryRepo.AssertExpectations(t)
}

func TestClockOut_NoOpenEntry(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	mockTimeEntryRepo.On("GetOpenEntry", uint(1)).Return(nil, errors.New("not found"))

	req := core.ClockOutRequest{EmployeeID: 1}
	entry, err := service.ClockOut(req)

	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Contains(t, err.Error(), "açık giriş kaydı bulunamadı")
}

func TestCreateTimeEntry_Normal(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	clockIn := time.Date(2026, 2, 26, 8, 0, 0, 0, time.UTC)
	clockOut := time.Date(2026, 2, 26, 16, 0, 0, 0, time.UTC)

	mockTimeEntryRepo.On("Create", mock.AnythingOfType("*core.TimeEntry")).Return(nil)

	req := core.TimeEntryRequest{
		EmployeeID: 1,
		ClockIn:    clockIn,
		ClockOut:   &clockOut,
		Notes:      "Normal mesai",
	}

	entry, err := service.CreateTimeEntry(req)

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "manual", entry.Source)
	assert.Equal(t, "pending", entry.Status)

	mockTimeEntryRepo.AssertExpectations(t)
}

func TestCreateTimeEntry_InvalidTimes(t *testing.T) {
	mockTimeEntryRepo := new(mocks.MockTimeEntryRepository)
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	service := services.NewTimekeepingService(mockTimeEntryRepo, mockScheduleRepo)

	clockIn := time.Date(2026, 2, 26, 16, 0, 0, 0, time.UTC)
	clockOut := time.Date(2026, 2, 26, 8, 0, 0, 0, time.UTC) // Before clock-in

	req := core.TimeEntryRequest{
		EmployeeID: 1,
		ClockIn:    clockIn,
		ClockOut:   &clockOut,
	}

	entry, err := service.CreateTimeEntry(req)

	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Contains(t, err.Error(), "çıkış saati giriş saatinden önce olamaz")
}

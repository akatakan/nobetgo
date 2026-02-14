package services_test

import (
	"testing"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	mocks "github.com/akatakan/nobetgo/tests/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGenerateSchedule(t *testing.T) {
	mockScheduleRepo := new(mocks.MockScheduleRepository)
	mockEmployeeRepo := new(mocks.MockEmployeeRepository)
	mockShiftRepo := new(mocks.MockShiftTypeRepository)

	service := services.NewSchedulerService(mockScheduleRepo, mockEmployeeRepo, mockShiftRepo)

	t.Run("Success", func(t *testing.T) {
		req := core.ScheduleRequest{Month: 1, Year: 2025}

		// Mock expected behaviors
		mockScheduleRepo.On("DeleteByMonthYear", 1, 2025).Return(nil)

		employees := []core.Employee{{Model: gorm.Model{ID: 1}, FirstName: "Dr. A"}, {Model: gorm.Model{ID: 2}, FirstName: "Dr. B"}}
		shiftTypes := []core.ShiftType{{Model: gorm.Model{ID: 1}, Name: "Nöbet"}}

		mockEmployeeRepo.On("List").Return(employees, nil)
		mockShiftRepo.On("List").Return(shiftTypes, nil)

		mockScheduleRepo.On("Create", mock.AnythingOfType("*core.Schedule")).Return(nil)

		schedules, err := service.GenerateSchedule(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, schedules)

		mockScheduleRepo.AssertCalled(t, "DeleteByMonthYear", 1, 2025)
		mockEmployeeRepo.AssertCalled(t, "List")
		mockShiftRepo.AssertCalled(t, "List")
	})
}

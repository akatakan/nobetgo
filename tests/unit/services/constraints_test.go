package services_test

import (
	"testing"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services/scheduler"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNoConsecutiveShifts(t *testing.T) {
	constraint := &scheduler.NoConsecutiveShifts{}

	emp := core.Employee{Model: gorm.Model{ID: 1}}
	date := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	// Case 1: Violation (Worked yesterday)
	schedule := []core.Schedule{
		{
			EmployeeID: 1,
			Date:       time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	penalty, isHard := constraint.CalculatePenalty(schedule, emp, date, core.ShiftType{})
	assert.True(t, isHard, "Should be a hard constraint violation")
	assert.Greater(t, penalty, 0.0)

	// Case 2: No Violation
	scheduleValid := []core.Schedule{
		{
			EmployeeID: 1,
			Date:       time.Date(2025, 1, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	penaltyValid, isHardValid := constraint.CalculatePenalty(scheduleValid, emp, date, core.ShiftType{})
	assert.False(t, isHardValid)
	assert.Equal(t, 0.0, penaltyValid)
}

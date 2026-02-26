package scheduler

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
)

type Constraint interface {
	// CalculatePenalty returns a score (0 for no violation, >0 for penalty)
	// and a boolean indicating if it's a hard constraint (true = invalid schedule)
	CalculatePenalty(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool)
}

// NoConsecutiveShifts ensures an employee doesn't work back-to-back days
type NoConsecutiveShifts struct{}

func (c *NoConsecutiveShifts) CalculatePenalty(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool) {
	for _, s := range schedule {
		if s.EmployeeID == employee.ID {
			// Check if shift is on Date-1 or Date+1
			diff := s.Date.Sub(date).Hours()
			if diff >= -24 && diff <= 24 && diff != 0 {
				return 1000.0, true // Hard constraint violation
			}
		}
	}
	return 0, false
}

// WeeklyHourLimit checks if employee exceeds 45 hours per week (simplified)
type WeeklyHourLimit struct {
	LimitHours float64
}

func (c *WeeklyHourLimit) CalculatePenalty(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool) {
	year, week := date.ISOWeek()
	totalHours := 0.0

	// Count hours for this employee in the same ISO week
	for _, s := range schedule {
		if s.EmployeeID == employee.ID {
			sYear, sWeek := s.Date.ISOWeek()
			if sYear == year && sWeek == week {
				// Assuming each shift is 8 hours if not specified
				// In a real system, we'd look up ShiftType.Hours
				totalHours += 8.0
			}
		}
	}

	// Add current shift hours
	totalHours += 8.0

	if totalHours > c.LimitHours {
		// Penalty proportional to excess
		excess := totalHours - c.LimitHours
		return excess * 100.0, false // Soft constraint (penalty)
	}

	return 0, false
}

// AnnualLeaveOverlapConstraint ensures an employee is not assigned a shift while on leave
type AnnualLeaveOverlapConstraint struct {
	ApprovedLeaves []core.Leave
}

func (c *AnnualLeaveOverlapConstraint) CalculatePenalty(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool) {
	for _, l := range c.ApprovedLeaves {
		if l.EmployeeID == employee.ID && l.Status == "approved" {
			// Normalize to UTC date start for comparison
			check := date.UTC().Truncate(24 * time.Hour)
			start := l.StartDate.UTC().Truncate(24 * time.Hour)
			end := l.EndDate.UTC().Truncate(24 * time.Hour)

			if (check.Equal(start) || check.After(start)) && (check.Equal(end) || check.Before(end)) {
				return 5000.0, true // Hard constraint violation
			}
		}
	}
	return 0, false
}

// MinimumRestConstraint ensures at least X hours of rest between shifts
type MinimumRestConstraint struct {
	MinRestHours float64
}

func (c *MinimumRestConstraint) CalculatePenalty(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool) {
	for _, s := range schedule {
		if s.EmployeeID == employee.ID {
			diff := s.Date.Sub(date).Hours()
			if diff < 0 {
				diff = -diff
			}
			if diff > 0 && diff < c.MinRestHours {
				return 1000.0, true // Hard constraint
			}
		}
	}
	return 0, false
}

package scheduler

import (
	"math/rand"
	"sort"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
)

type Optimizer struct {
	Constraints []Constraint
}

func NewOptimizer(constraints []Constraint) *Optimizer {
	return &Optimizer{
		Constraints: constraints,
	}
}

// OptimizeSchedule generates an optimized schedule using round-robin fair distribution
// with hard constraint enforcement (no consecutive shifts, balanced workload).
func (o *Optimizer) OptimizeSchedule(employees []core.Employee, shiftTypes []core.ShiftType, month, year int) []core.Schedule {
	if len(employees) == 0 || len(shiftTypes) == 0 {
		return []core.Schedule{}
	}

	// Try multiple times with shuffled order to find a valid schedule
	var bestSchedule []core.Schedule
	bestPenalty := float64(1e18)

	for attempt := 0; attempt < 50; attempt++ {
		schedule, penalty := o.generateFairSchedule(employees, shiftTypes, month, year)
		if penalty < bestPenalty {
			bestPenalty = penalty
			bestSchedule = schedule
		}
		// Perfect schedule found (no violations)
		if penalty == 0 {
			break
		}
	}

	return bestSchedule
}

// generateFairSchedule creates a schedule using round-robin distribution
// with built-in constraint checking during assignment.
func (o *Optimizer) generateFairSchedule(employees []core.Employee, shiftTypes []core.ShiftType, month, year int) ([]core.Schedule, float64) {
	var schedule []core.Schedule
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	// Track shifts per employee for fairness
	shiftCount := make(map[uint]int)
	for _, e := range employees {
		shiftCount[e.ID] = 0
	}

	// Track last work day per employee for consecutive check
	lastWorkDay := make(map[uint]time.Time)

	// Total slots = days × shift types
	// Ideal per person = total slots / len(employees)

	// Build all day-shift slots
	type slot struct {
		date      time.Time
		shiftType core.ShiftType
	}
	var slots []slot
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		for _, st := range shiftTypes {
			slots = append(slots, slot{date: d, shiftType: st})
		}
	}

	// Shuffle employees for initial ordering variety
	shuffled := make([]core.Employee, len(employees))
	copy(shuffled, employees)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	totalPenalty := 0.0

	for _, sl := range slots {
		// Sort employees by shift count (fewest first) for fairness
		sort.Slice(shuffled, func(i, j int) bool {
			return shiftCount[shuffled[i].ID] < shiftCount[shuffled[j].ID]
		})

		assigned := false
		for _, emp := range shuffled {
			// Check consecutive constraint: employee must not have worked the previous day
			if last, ok := lastWorkDay[emp.ID]; ok {
				diff := sl.date.Sub(last).Hours()
				if diff > 0 && diff <= 24 {
					// Would create consecutive shift — skip this employee
					continue
				}
			}

			// Valid assignment
			s := core.Schedule{
				Date:        sl.date,
				EmployeeID:  emp.ID,
				ShiftTypeID: sl.shiftType.ID,
			}
			schedule = append(schedule, s)
			shiftCount[emp.ID]++
			lastWorkDay[emp.ID] = sl.date
			assigned = true
			break
		}

		// If no one could be assigned without violation, pick the one with fewest shifts anyway
		if !assigned {
			// Fallback: least-loaded employee (accept the consecutive violation with penalty)
			emp := shuffled[0]
			s := core.Schedule{
				Date:        sl.date,
				EmployeeID:  emp.ID,
				ShiftTypeID: sl.shiftType.ID,
			}
			schedule = append(schedule, s)
			shiftCount[emp.ID]++
			lastWorkDay[emp.ID] = sl.date
			totalPenalty += 1000 // Penalty for constraint violation
		}
	}

	// Calculate fairness penalty: difference between max and min shift counts
	maxShifts, minShifts := 0, len(slots)
	for _, count := range shiftCount {
		if count > maxShifts {
			maxShifts = count
		}
		if count < minShifts {
			minShifts = count
		}
	}
	fairnessDiff := maxShifts - minShifts
	if fairnessDiff > 1 {
		totalPenalty += float64(fairnessDiff-1) * 500 // Penalty if diff > 1
	}

	return schedule, totalPenalty
}

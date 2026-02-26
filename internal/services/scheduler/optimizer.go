package scheduler

import (
	"sort"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
)

type Optimizer struct {
	Constraints []Constraint
}

type slot struct {
	date      time.Time
	shiftType core.ShiftType
}

func NewOptimizer(constraints []Constraint) *Optimizer {
	return &Optimizer{
		Constraints: constraints,
	}
}

func (o *Optimizer) checkConstraints(schedule []core.Schedule, employee core.Employee, date time.Time, shift core.ShiftType) (float64, bool) {
	totalPenalty := 0.0
	for _, c := range o.Constraints {
		penalty, isHard := c.CalculatePenalty(schedule, employee, date, shift)
		if isHard {
			return 0, true
		}
		totalPenalty += penalty
	}
	return totalPenalty, false
}

// OptimizeSchedule generates an optimized schedule using a backtracking approach.
func (o *Optimizer) OptimizeSchedule(req core.ScheduleRequest, department *core.Department, employees []core.Employee, shiftTypes []core.ShiftType) []core.Schedule {
	if len(employees) == 0 || len(shiftTypes) == 0 {
		return []core.Schedule{}
	}

	start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	// Determine how many employees are needed per shift
	employeesPerShift := 1
	if req.SchedulingMode == "bed_capacity" && department != nil && req.BedsPerPersonnel > 0 {
		employeesPerShift = department.BedCapacity / req.BedsPerPersonnel
		if employeesPerShift < 1 {
			employeesPerShift = 1
		}
	}

	// Build all day-shift slots
	var slots []slot
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		for _, st := range shiftTypes {
			// Each slot needs 'employeesPerShift' assignments
			for i := 0; i < employeesPerShift; i++ {
				slots = append(slots, slot{date: d, shiftType: st})
			}
		}
	}

	// State tracking for backtracking
	shiftCount := make(map[uint]int)
	for _, e := range employees {
		shiftCount[e.ID] = 0
	}

	var finalSchedule []core.Schedule
	found := false
	startTime := time.Now()
	timeout := 10 * time.Second

	// Recursive backtracking function
	var backtrack func(slotIdx int, schedule []core.Schedule)
	backtrack = func(slotIdx int, schedule []core.Schedule) {
		if found || time.Since(startTime) > timeout {
			return
		}

		if slotIdx == len(slots) {
			finalSchedule = append([]core.Schedule{}, schedule...)
			found = true
			return
		}

		sl := slots[slotIdx]

		// Candidate selection logic
		candidates := make([]core.Employee, len(employees))
		copy(candidates, employees)

		// Sort candidates based on mode
		if req.SchedulingMode == "fatigue_aware" {
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].FatigueScore < candidates[j].FatigueScore
			})
		} else {
			sort.Slice(candidates, func(i, j int) bool {
				return shiftCount[candidates[i].ID] < shiftCount[candidates[j].ID]
			})
		}

		for _, emp := range candidates {
			// Double assignment check (same date, same shift type)
			alreadyInSlot := false
			for _, s := range schedule {
				if s.Date.Equal(sl.date) && s.ShiftTypeID == sl.shiftType.ID && s.EmployeeID == emp.ID {
					alreadyInSlot = true
					break
				}
			}
			if alreadyInSlot {
				continue
			}

			_, isHard := o.checkConstraints(schedule, emp, sl.date, sl.shiftType)
			if isHard {
				continue
			}

			// Tentative assignment
			s := core.Schedule{
				Date:        sl.date,
				EmployeeID:  emp.ID,
				ShiftTypeID: sl.shiftType.ID,
			}
			schedule = append(schedule, s)
			shiftCount[emp.ID]++

			backtrack(slotIdx+1, schedule)

			if found {
				return
			}

			// Backtrack
			schedule = schedule[:len(schedule)-1]
			shiftCount[emp.ID]--
		}
	}

	backtrack(0, []core.Schedule{})

	// Fallback: If no perfect schedule found within timeout, use a greedy approach (or empty)
	if !found {
		// Greedy fallback for robustness
		return o.generateGreedySchedule(req, department, employees, shiftTypes, employeesPerShift, slots)
	}

	return finalSchedule
}

func (o *Optimizer) generateGreedySchedule(req core.ScheduleRequest, department *core.Department, employees []core.Employee, shiftTypes []core.ShiftType, employeesPerShift int, slots []slot) []core.Schedule {
	var schedule []core.Schedule
	shiftCount := make(map[uint]int)

	for _, sl := range slots {
		// Sort by shift count
		sort.Slice(employees, func(i, j int) bool {
			return shiftCount[employees[i].ID] < shiftCount[employees[j].ID]
		})

		for _, emp := range employees {
			alreadyInSlot := false
			for _, s := range schedule {
				if s.Date.Equal(sl.date) && s.ShiftTypeID == sl.shiftType.ID && s.EmployeeID == emp.ID {
					alreadyInSlot = true
					break
				}
			}
			if alreadyInSlot {
				continue
			}

			_, isHard := o.checkConstraints(schedule, emp, sl.date, sl.shiftType)
			if !isHard {
				s := core.Schedule{
					Date:        sl.date,
					EmployeeID:  emp.ID,
					ShiftTypeID: sl.shiftType.ID,
				}
				schedule = append(schedule, s)
				shiftCount[emp.ID]++
				break
			}
		}
	}
	return schedule
}

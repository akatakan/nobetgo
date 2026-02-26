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
func (o *Optimizer) OptimizeSchedule(req core.ScheduleRequest, department *core.Department, employees []core.Employee, shiftTypes []core.ShiftType) []core.Schedule {
	if len(employees) == 0 || len(shiftTypes) == 0 {
		return []core.Schedule{}
	}

	// Try multiple times with shuffled order to find a valid schedule
	var bestSchedule []core.Schedule
	bestPenalty := float64(1e18)

	for attempt := 0; attempt < 50; attempt++ {
		var schedule []core.Schedule
		var penalty float64

		if req.SchedulingMode == "fatigue_aware" {
			schedule, penalty = o.generateFatigueAwareSchedule(req, department, employees, shiftTypes)
		} else {
			schedule, penalty = o.generateFairSchedule(req, department, employees, shiftTypes)
		}

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
func (o *Optimizer) generateFairSchedule(req core.ScheduleRequest, department *core.Department, employees []core.Employee, shiftTypes []core.ShiftType) ([]core.Schedule, float64) {
	var schedule []core.Schedule
	start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	// Determine how many employees are needed per shift
	employeesPerShift := 1
	if req.SchedulingMode == "bed_capacity" && department != nil && req.BedsPerPersonnel > 0 {
		employeesPerShift = department.BedCapacity / req.BedsPerPersonnel
		if employeesPerShift < 1 {
			employeesPerShift = 1 // Safety fallback
		}
	}

	// Track shifts per employee for fairness
	shiftCount := make(map[uint]int)
	for _, e := range employees {
		shiftCount[e.ID] = 0
	}

	// Track last work day per employee for consecutive check
	lastWorkDay := make(map[uint]time.Time)

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
		assignedCount := 0

		// Repeat assignment until we fill the required number of employees for this slot
		for assignedCount < employeesPerShift {
			// Sort employees by shift count (fewest first) for fairness
			sort.Slice(shuffled, func(i, j int) bool {
				return shiftCount[shuffled[i].ID] < shiftCount[shuffled[j].ID]
			})

			assignedThisRound := false
			for _, emp := range shuffled {
				// Prevent assigning the same employee twice to the SAME slot
				alreadyAssignedToThisSlot := false
				for _, existing := range schedule {
					if existing.Date.Equal(sl.date) && existing.ShiftTypeID == sl.shiftType.ID && existing.EmployeeID == emp.ID {
						alreadyAssignedToThisSlot = true
						break
					}
				}
				if alreadyAssignedToThisSlot {
					continue
				}

				// Check consecutive constraint: employee must not have worked the previous day
				if last, ok := lastWorkDay[emp.ID]; ok {
					diff := sl.date.Sub(last).Hours()
					// We only enforce 24h rest rule if diff > 0 (it's not the same day slot)
					if diff > 0 && diff <= 24 {
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
				assignedThisRound = true
				assignedCount++
				break
			}

			// If no one could be assigned without violation, pick the one with fewest shifts anyway
			if !assignedThisRound {
				// Find an employee not already in this slot
				var fallbackEmp *core.Employee
				for _, emp := range shuffled {
					alreadyAssigned := false
					for _, existing := range schedule {
						if existing.Date.Equal(sl.date) && existing.ShiftTypeID == sl.shiftType.ID && existing.EmployeeID == emp.ID {
							alreadyAssigned = true
							break
						}
					}
					if !alreadyAssigned {
						fallbackEmp = &emp
						break
					}
				}

				if fallbackEmp != nil {
					s := core.Schedule{
						Date:        sl.date,
						EmployeeID:  fallbackEmp.ID,
						ShiftTypeID: sl.shiftType.ID,
					}
					schedule = append(schedule, s)
					shiftCount[fallbackEmp.ID]++
					lastWorkDay[fallbackEmp.ID] = sl.date
					totalPenalty += 1000 // Penalty for constraint violation
					assignedCount++
				} else {
					// Extremely constrained (e.g. required per shift > total employees)
					// Break inner loop to avoid infinite loop
					break
				}
			}
		}
	}

	// Calculate fairness penalty: difference between max and min shift counts
	maxShifts, minShifts := 0, len(slots)*employeesPerShift
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

// generateFatigueAwareSchedule uses "Fatigue Battery" concepts.
// It assigns penalty based on current fatigue, avoids assigning if fatigue > 50,
// and grants a HeroPoint if it's struggling to find an available person.
func (o *Optimizer) generateFatigueAwareSchedule(req core.ScheduleRequest, department *core.Department, employees []core.Employee, shiftTypes []core.ShiftType) ([]core.Schedule, float64) {
	var schedule []core.Schedule
	start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	employeesPerShift := 1
	if req.BedsPerPersonnel > 0 && department != nil {
		employeesPerShift = department.BedCapacity / req.BedsPerPersonnel
		if employeesPerShift < 1 {
			employeesPerShift = 1
		}
	}

	// Copy initial fatigue & hero points so we can mutate safely per generated variant
	fatigue := make(map[uint]int)
	hero := make(map[uint]int)
	for _, e := range employees {
		fatigue[e.ID] = e.FatigueScore
		hero[e.ID] = e.HeroPoint
	}

	shiftCount := make(map[uint]int)
	lastWorkDay := make(map[uint]time.Time)

	type slot struct {
		date      time.Time
		shiftType core.ShiftType
	}
	var slots []slot
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		for _, st := range shiftTypes {
			slots = append(slots, slot{date: d, shiftType: st})
		}
		// Apply -2 rest score if employee had no shift that day
		for _, e := range employees {
			fatigue[e.ID] -= 2
			if fatigue[e.ID] < 0 {
				fatigue[e.ID] = 0
			}
		}
	}

	shuffled := make([]core.Employee, len(employees))
	copy(shuffled, employees)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	totalPenalty := 0.0

	for _, sl := range slots {
		assignedCount := 0
		for assignedCount < employeesPerShift {

			// Soft Constraint logic: Sort by lowest dynamic penalty
			// Penalty = (Fatigue_Current * Shift Difficulty) - (HeroPoint * 10)
			sort.Slice(shuffled, func(i, j int) bool {
				diff := 1 // Default normal shift difficulty
				if sl.shiftType.IsNightShift {
					diff = 3
				} else if sl.date.Weekday() == time.Saturday || sl.date.Weekday() == time.Sunday {
					diff = 2
				}

				penaltyI := (fatigue[shuffled[i].ID] * diff) - (hero[shuffled[i].ID] * 10)
				penaltyJ := (fatigue[shuffled[j].ID] * diff) - (hero[shuffled[j].ID] * 10)
				return penaltyI < penaltyJ
			})

			assignedThisRound := false
			for _, emp := range shuffled {
				// Prevent double assignment
				alreadyAssigned := false
				for _, existing := range schedule {
					if existing.Date.Equal(sl.date) && existing.ShiftTypeID == sl.shiftType.ID && existing.EmployeeID == emp.ID {
						alreadyAssigned = true
						break
					}
				}
				if alreadyAssigned {
					continue
				}

				// Hard Rest Constraint: 24 hours between shifts
				if last, ok := lastWorkDay[emp.ID]; ok {
					diffHours := sl.date.Sub(last).Hours()
					if diffHours > 0 && diffHours <= 24 {
						continue
					}
				}

				// Hard Fatigue Constraint: > 50 means burnout risk
				diff := 1
				if sl.shiftType.IsNightShift {
					diff = 3
				} else if sl.date.Weekday() == time.Saturday || sl.date.Weekday() == time.Sunday {
					diff = 2
				}

				if fatigue[emp.ID]+diff > 50 {
					continue // Riskli, zorunlu dinlenme
				}

				// Valid Assignment
				s := core.Schedule{
					Date:        sl.date,
					EmployeeID:  emp.ID,
					ShiftTypeID: sl.shiftType.ID,
				}
				schedule = append(schedule, s)
				shiftCount[emp.ID]++
				lastWorkDay[emp.ID] = sl.date

				// Update fatigue
				fatigue[emp.ID] += diff

				// Deduct hero point usage if they had any
				if hero[emp.ID] > 0 {
					hero[emp.ID]--
				}

				totalPenalty += float64((fatigue[emp.ID] * diff) - (hero[emp.ID] * 10))

				assignedThisRound = true
				assignedCount++
				break
			}

			if !assignedThisRound {
				// Hero System Substitution
				// Everyone was too tired or resting. Find someone valid (even if fatigue > 50).
				var fallbackEmp *core.Employee
				for _, emp := range shuffled {
					alreadyAssigned := false
					for _, existing := range schedule {
						if existing.Date.Equal(sl.date) && existing.ShiftTypeID == sl.shiftType.ID && existing.EmployeeID == emp.ID {
							alreadyAssigned = true
							break
						}
					}
					// Only basic rest rule applies
					if last, ok := lastWorkDay[emp.ID]; ok {
						diffHours := sl.date.Sub(last).Hours()
						if diffHours > 0 && diffHours <= 24 {
							continue
						}
					}
					if !alreadyAssigned {
						fallbackEmp = &emp
						break
					}
				}

				if fallbackEmp != nil {
					s := core.Schedule{
						Date:        sl.date,
						EmployeeID:  fallbackEmp.ID,
						ShiftTypeID: sl.shiftType.ID,
					}
					schedule = append(schedule, s)
					shiftCount[fallbackEmp.ID]++
					lastWorkDay[fallbackEmp.ID] = sl.date

					// Grand them a hero point for substituting a hard shift
					hero[fallbackEmp.ID]++
					fatigue[fallbackEmp.ID] += 5 // Heavy toll for hero shift

					totalPenalty += 2000 // Huge penalty for fallback
					assignedCount++
				} else {
					break
				}
			}
		}
	}

	return schedule, totalPenalty
}

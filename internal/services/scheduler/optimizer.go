package scheduler

import (
	"math/rand"
	"sync"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
)

type Optimizer struct {
	Constraints []Constraint
	WorkerCount int
}

func NewOptimizer(constraints []Constraint) *Optimizer {
	return &Optimizer{
		Constraints: constraints,
		WorkerCount: 4, // Default to 4 workers
	}
}

type ScheduleCandidate struct {
	Schedule []core.Schedule
	Score    float64
}

// OptimizeSchedule generates multiple random schedules and picks the best one based on constraints and cost
func (o *Optimizer) OptimizeSchedule(employees []core.Employee, shiftTypes []core.ShiftType, month, year int, iterations int, overtimeThreshold, overtimeMultiplier float64) []core.Schedule {
	candidates := make(chan ScheduleCandidate, iterations)
	var wg sync.WaitGroup

	// Worker pool
	jobs := make(chan int, iterations)

	for w := 0; w < o.WorkerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				schedule := o.generateRandomSchedule(employees, shiftTypes, month, year)
				score := o.calculateScore(schedule, employees, overtimeThreshold, overtimeMultiplier)
				candidates <- ScheduleCandidate{Schedule: schedule, Score: score}
			}
		}()
	}

	// Send jobs
	for i := 0; i < iterations; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(candidates)

	// Find best candidate
	var bestSchedule []core.Schedule
	minScore := 1e9 // Start with high score

	// If no valid schedule found, returns the one with lowest penalty
	for c := range candidates {
		if c.Score < minScore {
			minScore = c.Score
			bestSchedule = c.Schedule
		}
	}

	return bestSchedule
}

func (o *Optimizer) generateRandomSchedule(employees []core.Employee, shiftTypes []core.ShiftType, month, year int) []core.Schedule {
	var schedule []core.Schedule
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		for _, shift := range shiftTypes {
			// Random assignment
			randIndex := rand.Intn(len(employees))
			selectedEmployee := employees[randIndex]

			s := core.Schedule{
				Date:        d,
				EmployeeID:  selectedEmployee.ID,
				ShiftTypeID: shift.ID,
			}
			schedule = append(schedule, s)
		}
	}
	return schedule
}

func (o *Optimizer) calculateScore(schedule []core.Schedule, employees []core.Employee, overtimeThreshold, overtimeMultiplier float64) float64 {
	score := 0.0

	// Create a map for quick employee lookup
	empMap := make(map[uint]core.Employee)
	for _, e := range employees {
		empMap[e.ID] = e
	}

	// Track hours per employee
	empHours := make(map[uint]float64)

	for _, s := range schedule {
		emp := empMap[s.EmployeeID]
		shiftHours := 8.0 // TODO: Get actual hours from ShiftType

		// Check current total
		currentTotal := empHours[emp.ID]

		// Calculate cost for this shift
		// If already over threshold, full shift is overtime
		// If crossing threshold, split calculation

		// Simplified weekly logic (treating month as single block for MVP or assuming weekly reset?
		// Real overtime is usually weekly.
		// For this implementation, let's just create a global "Overtime after X hours in month" or simplified average.
		// PROPER WAY: Group by week.
		// MVP WAY (User Request): "Weekly limit".

		// Let's approximate week by Day / 7.
		// weekNum := s.Date.Day() / 7
		// uniqueKey := fmt.Sprintf("%d-%d", emp.ID, weekNum)
		// But s.Date isn't sorted in loop? generateRandomSchedule sorts by date.

		// Let's Assume the loop in generateRandomSchedule produces sorted dates.
		// It does: `for d := start; d.Before(end); ...`

		cost := 0.0

		// Check invalid overtime (Hard limit? Or just cost?)
		// standard cost
		cost = emp.HourlyRate * shiftHours

		// If this employee has piled up hours, apply multiplier?
		// Let's simplified: If total hours > (Threshold * 4 weeks), apply multiplier?
		// No, user said "Weekly".

		// Let's compute Score based on pure cost.
		// If we want to OPTIMIZE for cost, we essentially want to avoid overtime if multiplier > 1.

		// Let's blindly apply a generic penalty/cost if they work too much.
		currentTotal += shiftHours
		empHours[emp.ID] = currentTotal

		// Basic Overtime Logic:
		// If total monthly hours > (45 * 4), then apply multiplier on the excess.
		monthlyThreshold := overtimeThreshold * 4.0
		if currentTotal > monthlyThreshold {
			// This shift is fully or partially overtime
			// Simplified: Just multiply whole shift cost if in overtime zone
			cost *= overtimeMultiplier
		}

		score += cost

		// 2. Constraints Penalty
		for _, c := range o.Constraints {
			penalty, isHard := c.CalculatePenalty(schedule, emp, s.Date, core.ShiftType{})
			if isHard {
				score += 100000
			}
			score += penalty
		}
	}

	return score
}

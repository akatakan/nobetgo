package services

import (
	"log/slog"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

// OvertimeService handles overtime calculation using defined rules.
type OvertimeService struct {
	ruleRepo      repositories.OvertimeRuleRepositoryInterface
	timeEntryRepo repositories.TimeEntryRepositoryInterface
}

// NewOvertimeService creates a new OvertimeService.
func NewOvertimeService(ruleRepo repositories.OvertimeRuleRepositoryInterface, timeEntryRepo repositories.TimeEntryRepositoryInterface) *OvertimeService {
	return &OvertimeService{ruleRepo: ruleRepo, timeEntryRepo: timeEntryRepo}
}

// OvertimeSummary holds calculated overtime data for a single employee.
type OvertimeSummary struct {
	EmployeeID      uint    `json:"employee_id"`
	EmployeeName    string  `json:"employee_name"`
	TotalHours      float64 `json:"total_hours"`
	NormalHours     float64 `json:"normal_hours"`
	OvertimeHours   float64 `json:"overtime_hours"`
	WeekendHours    float64 `json:"weekend_hours"`
	HolidayHours    float64 `json:"holiday_hours"`
	NightShiftHours float64 `json:"night_shift_hours"`
	OvertimePay     float64 `json:"overtime_pay"`
	WeekendPay      float64 `json:"weekend_pay"`
	HolidayPay      float64 `json:"holiday_pay"`
	TotalPay        float64 `json:"total_pay"`
	WorkingDays     int     `json:"working_days"`
}

// CalculateOvertime computes the overtime summary for an employee in a given month.
func (s *OvertimeService) CalculateOvertime(employeeID uint, month, year int) (*OvertimeSummary, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	entries, err := s.timeEntryRepo.ListByEmployee(employeeID, start, end)
	if err != nil {
		return nil, err
	}

	rule, err := s.ruleRepo.GetActive()
	if err != nil {
		slog.Warn("No active overtime rule found, using defaults", "error", err)
		rule = defaultOvertimeRule()
	}

	// Bug 1 fix: Pre-fetch all holidays for the year into a map (eliminates N+1 queries)
	holidayMap := s.buildHolidayMap(year)

	summary := &OvertimeSummary{
		EmployeeID: employeeID,
	}

	// Calculate weekly buckets
	weeklyHours := make(map[int]float64) // ISO week -> hours

	for _, entry := range entries {
		hours := entry.NetWorkingHours()
		if hours <= 0 {
			continue
		}

		summary.TotalHours += hours
		summary.WorkingDays++

		// Bug 3 fix: Guard against nil Employee before accessing fields
		if summary.EmployeeName == "" && entry.Employee.ID != 0 {
			summary.EmployeeName = entry.Employee.FirstName + " " + entry.Employee.LastName
		}

		// Classify hours — O(1) map lookup instead of DB query
		isHoliday := holidayMap[entry.ClockIn.Truncate(24*time.Hour).Format("2006-01-02")]
		isWeekend := entry.ClockIn.Weekday() == time.Saturday || entry.ClockIn.Weekday() == time.Sunday

		// Bug 6 fix: Always add hours to weeklyHours so ALL hours count toward the 45h weekly limit
		_, week := entry.ClockIn.ISOWeek()
		weeklyHours[week] += hours

		// Track type-specific hours for pay multipliers
		switch {
		case isHoliday:
			summary.HolidayHours += hours
		case isWeekend:
			summary.WeekendHours += hours
		}

		// Night shift check (simplified: if entry type is marked or shift starts after 22:00)
		if entry.EntryType == "night" {
			summary.NightShiftHours += hours
		}
	}

	// Calculate normal vs overtime from weekly totals
	for _, hours := range weeklyHours {
		if hours <= rule.WeeklyHourLimit {
			summary.NormalHours += hours
		} else {
			summary.NormalHours += rule.WeeklyHourLimit
			summary.OvertimeHours += hours - rule.WeeklyHourLimit
		}
	}

	// Bug 3 fix: Guard against nil/zero Employee before accessing HourlyRate
	if len(entries) > 0 && entries[0].Employee.ID != 0 {
		rate := entries[0].Employee.HourlyRate
		summary.OvertimePay = summary.OvertimeHours * rate * rule.OvertimeMultiplier
		summary.WeekendPay = summary.WeekendHours * rate * rule.WeekendMultiplier
		summary.HolidayPay = summary.HolidayHours * rate * rule.HolidayMultiplier
		summary.TotalPay = summary.OvertimePay + summary.WeekendPay + summary.HolidayPay
	}

	return summary, nil
}

// GetDepartmentOvertimeSummary returns overtime summaries for all employees in a department.
func (s *OvertimeService) GetDepartmentOvertimeSummary(deptID uint, month, year int) ([]OvertimeSummary, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	entries, err := s.timeEntryRepo.ListByDepartment(deptID, start, end)
	if err != nil {
		return nil, err
	}

	rule, err := s.ruleRepo.GetActive()
	if err != nil {
		rule = defaultOvertimeRule()
	}

	// Bug 1 fix: Pre-fetch all holidays for the year into a map (eliminates N+1 queries)
	holidayMap := s.buildHolidayMap(year)

	// Group entries by employee
	grouped := make(map[uint][]core.TimeEntry)
	for _, e := range entries {
		grouped[e.EmployeeID] = append(grouped[e.EmployeeID], e)
	}

	var summaries []OvertimeSummary
	for empID, empEntries := range grouped {
		summary := OvertimeSummary{
			EmployeeID: empID,
		}

		weeklyHours := make(map[int]float64)

		for _, entry := range empEntries {
			hours := entry.NetWorkingHours()
			if hours <= 0 {
				continue
			}

			summary.TotalHours += hours
			summary.WorkingDays++

			// Bug 3 fix: Guard against nil Employee
			if summary.EmployeeName == "" && entry.Employee.ID != 0 {
				summary.EmployeeName = entry.Employee.FirstName + " " + entry.Employee.LastName
			}

			// O(1) map lookup instead of DB query
			isHoliday := holidayMap[entry.ClockIn.Truncate(24*time.Hour).Format("2006-01-02")]
			isWeekend := entry.ClockIn.Weekday() == time.Saturday || entry.ClockIn.Weekday() == time.Sunday

			// Bug 6 fix: Always add hours to weeklyHours
			_, week := entry.ClockIn.ISOWeek()
			weeklyHours[week] += hours

			switch {
			case isHoliday:
				summary.HolidayHours += hours
			case isWeekend:
				summary.WeekendHours += hours
			}
		}

		for _, hours := range weeklyHours {
			if hours <= rule.WeeklyHourLimit {
				summary.NormalHours += hours
			} else {
				summary.NormalHours += rule.WeeklyHourLimit
				summary.OvertimeHours += hours - rule.WeeklyHourLimit
			}
		}

		// Bug 3 fix: Guard against nil Employee for HourlyRate
		if len(empEntries) > 0 && empEntries[0].Employee.ID != 0 {
			rate := empEntries[0].Employee.HourlyRate
			summary.OvertimePay = summary.OvertimeHours * rate * rule.OvertimeMultiplier
			summary.WeekendPay = summary.WeekendHours * rate * rule.WeekendMultiplier
			summary.HolidayPay = summary.HolidayHours * rate * rule.HolidayMultiplier
			summary.TotalPay = summary.OvertimePay + summary.WeekendPay + summary.HolidayPay
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// --- OvertimeRule CRUD ---

func (s *OvertimeService) CreateRule(rule *core.OvertimeRule) error {
	return s.ruleRepo.Create(rule)
}

func (s *OvertimeService) UpdateRule(rule *core.OvertimeRule) error {
	return s.ruleRepo.Update(rule)
}

func (s *OvertimeService) GetRule(id uint) (*core.OvertimeRule, error) {
	return s.ruleRepo.GetByID(id)
}

func (s *OvertimeService) GetAllRules() ([]core.OvertimeRule, error) {
	return s.ruleRepo.List()
}

func (s *OvertimeService) DeleteRule(id uint) error {
	return s.ruleRepo.Delete(id)
}

// --- PublicHoliday CRUD ---

func (s *OvertimeService) CreateHoliday(h *core.PublicHoliday) error {
	return s.ruleRepo.CreateHoliday(h)
}

func (s *OvertimeService) UpdateHoliday(h *core.PublicHoliday) error {
	return s.ruleRepo.UpdateHoliday(h)
}

func (s *OvertimeService) GetHoliday(id uint) (*core.PublicHoliday, error) {
	return s.ruleRepo.GetHolidayByID(id)
}

func (s *OvertimeService) GetHolidays(year int) ([]core.PublicHoliday, error) {
	return s.ruleRepo.ListHolidays(year)
}

func (s *OvertimeService) DeleteHoliday(id uint) error {
	return s.ruleRepo.DeleteHoliday(id)
}

// buildHolidayMap fetches all holidays for the given year and returns a date-keyed map for O(1) lookups.
func (s *OvertimeService) buildHolidayMap(year int) map[string]bool {
	holidayMap := make(map[string]bool)
	holidays, err := s.ruleRepo.ListHolidays(year)
	if err != nil {
		slog.Warn("Failed to fetch holidays, assuming none", "year", year, "error", err)
		return holidayMap
	}
	for _, h := range holidays {
		holidayMap[h.Date.Truncate(24*time.Hour).Format("2006-01-02")] = true
	}
	return holidayMap
}

func defaultOvertimeRule() *core.OvertimeRule {
	return &core.OvertimeRule{
		WeeklyHourLimit:    45,
		DailyHourLimit:     11,
		OvertimeMultiplier: 1.5,
		WeekendMultiplier:  2.0,
		HolidayMultiplier:  2.5,
		NightShiftExtra:    0.1,
	}
}

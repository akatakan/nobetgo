package services

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

// ReportingService provides analytics and reports for work hours, absences, and trends.
type ReportingService struct {
	timeEntryRepo repositories.TimeEntryRepositoryInterface
	leaveRepo     repositories.LeaveRepositoryInterface
	overtimeRepo  repositories.OvertimeRuleRepositoryInterface
}

// NewReportingService creates a new ReportingService.
func NewReportingService(
	timeEntryRepo repositories.TimeEntryRepositoryInterface,
	leaveRepo repositories.LeaveRepositoryInterface,
	overtimeRepo repositories.OvertimeRuleRepositoryInterface,
) *ReportingService {
	return &ReportingService{
		timeEntryRepo: timeEntryRepo,
		leaveRepo:     leaveRepo,
		overtimeRepo:  overtimeRepo,
	}
}

// WorkHoursReport holds aggregated work hours data.
type WorkHoursReport struct {
	Month       int                  `json:"month"`
	Year        int                  `json:"year"`
	Employees   []EmployeeWorkReport `json:"employees"`
	TotalHours  float64              `json:"total_hours"`
	WorkingDays int                  `json:"working_days"`
}

// EmployeeWorkReport holds work hours data for a single employee.
type EmployeeWorkReport struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Department   string  `json:"department"`
	TotalHours   float64 `json:"total_hours"`
	WorkingDays  int     `json:"working_days"`
	AvgDaily     float64 `json:"avg_daily_hours"`
}

// CostAnalysisReport holds cost data based on working hours.
type CostAnalysisReport struct {
	Month      int                  `json:"month"`
	Year       int                  `json:"year"`
	Employees  []EmployeeCostDetail `json:"employees"`
	TotalCost  float64              `json:"total_cost"`
	TotalHours float64              `json:"total_hours"`
}

// EmployeeCostDetail holds cost details for a single employee.
type EmployeeCostDetail struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Department   string  `json:"department"`
	TotalHours   float64 `json:"total_hours"`
	HourlyRate   float64 `json:"hourly_rate"`
	TotalCost    float64 `json:"total_cost"`
}

// AbsenceReport holds absence data by department.
type AbsenceReport struct {
	Month     int                     `json:"month"`
	Year      int                     `json:"year"`
	Employees []EmployeeAbsenceReport `json:"employees"`
	TotalDays float64                 `json:"total_days"`
}

// EmployeeAbsenceReport holds absence data for a single employee.
type EmployeeAbsenceReport struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	LeaveType    string  `json:"leave_type"`
	TotalDays    float64 `json:"total_days"`
}

// EmployeeSummaryReport holds a summary for a single employee.
type EmployeeSummaryReport struct {
	EmployeeID    uint    `json:"employee_id"`
	EmployeeName  string  `json:"employee_name"`
	Department    string  `json:"department"`
	TotalHours    float64 `json:"total_hours"`
	OvertimeHours float64 `json:"overtime_hours"`
	LeaveDays     float64 `json:"leave_days"`
	WorkingDays   int     `json:"working_days"`
}

// TrendData holds monthly trend data for analysis.
type TrendData struct {
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	TotalHours    float64 `json:"total_hours"`
	OvertimeHours float64 `json:"overtime_hours"`
	AbsenceDays   float64 `json:"absence_days"`
	WorkingDays   int     `json:"working_days"`
}

// GetWorkHoursReport generates a work hours report for a department in a given month.
func (s *ReportingService) GetWorkHoursReport(deptID uint, month, year int) (*WorkHoursReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var entries []core.TimeEntry
	var err error

	if deptID > 0 {
		entries, err = s.timeEntryRepo.ListByDepartment(deptID, start, end)
	} else {
		entries, err = s.timeEntryRepo.ListByDateRange(start, end)
	}
	if err != nil {
		return nil, err
	}

	report := &WorkHoursReport{Month: month, Year: year}

	// Group by employee
	grouped := make(map[uint]*EmployeeWorkReport)
	for _, e := range entries {
		hours := e.NetWorkingHours()
		emp, ok := grouped[e.EmployeeID]
		if !ok {
			emp = &EmployeeWorkReport{
				EmployeeID:   e.EmployeeID,
				EmployeeName: e.Employee.FirstName + " " + e.Employee.LastName,
				Department:   e.Employee.Department.Name,
			}
			grouped[e.EmployeeID] = emp
		}
		emp.TotalHours += hours
		emp.WorkingDays++
		report.TotalHours += hours
		report.WorkingDays++
	}

	for _, emp := range grouped {
		if emp.WorkingDays > 0 {
			emp.AvgDaily = emp.TotalHours / float64(emp.WorkingDays)
		}
		report.Employees = append(report.Employees, *emp)
	}

	return report, nil
}

// GetCostAnalysis calculates the monetary cost of employee hours for a given month.
func (s *ReportingService) GetCostAnalysis(deptID uint, month, year int) (*CostAnalysisReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var entries []core.TimeEntry
	var err error

	if deptID > 0 {
		entries, err = s.timeEntryRepo.ListByDepartment(deptID, start, end)
	} else {
		entries, err = s.timeEntryRepo.ListByDateRange(start, end)
	}
	if err != nil {
		return nil, err
	}

	report := &CostAnalysisReport{Month: month, Year: year}
	grouped := make(map[uint]*EmployeeCostDetail)

	for _, e := range entries {
		hours := e.NetWorkingHours()
		emp, ok := grouped[e.EmployeeID]
		if !ok {
			emp = &EmployeeCostDetail{
				EmployeeID:   e.EmployeeID,
				EmployeeName: e.Employee.FirstName + " " + e.Employee.LastName,
				Department:   e.Employee.Department.Name,
				HourlyRate:   e.Employee.HourlyRate,
			}
			grouped[e.EmployeeID] = emp
		}
		emp.TotalHours += hours
		report.TotalHours += hours
	}

	for _, emp := range grouped {
		emp.TotalCost = emp.TotalHours * emp.HourlyRate
		report.TotalCost += emp.TotalCost
		report.Employees = append(report.Employees, *emp)
	}

	return report, nil
}

// GetAbsenceReport generates an absence report for a department in a given month.
func (s *ReportingService) GetAbsenceReport(deptID uint, month, year int) (*AbsenceReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	var leaves []core.Leave
	var err error

	if deptID > 0 {
		leaves, err = s.leaveRepo.ListByDepartment(deptID, start, end)
	} else {
		leaves, err = s.leaveRepo.ListByEmployee(0, start, end) // fallback: empty for all
	}
	if err != nil {
		return nil, err
	}

	report := &AbsenceReport{Month: month, Year: year}

	for _, l := range leaves {
		if l.Status != "approved" {
			continue
		}
		report.Employees = append(report.Employees, EmployeeAbsenceReport{
			EmployeeID:   l.EmployeeID,
			EmployeeName: l.Employee.FirstName + " " + l.Employee.LastName,
			LeaveType:    l.LeaveType.Name,
			TotalDays:    l.TotalDays,
		})
		report.TotalDays += l.TotalDays
	}

	return report, nil
}

// GetEmployeeSummary generates a comprehensive summary for a single employee.
func (s *ReportingService) GetEmployeeSummary(employeeID uint, month, year int) (*EmployeeSummaryReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	entries, err := s.timeEntryRepo.ListByEmployee(employeeID, start, end)
	if err != nil {
		return nil, err
	}

	leaves, err := s.leaveRepo.ListByEmployee(employeeID, start, end)
	if err != nil {
		return nil, err
	}

	rule, err := s.overtimeRepo.GetActive()
	if err != nil {
		rule = &core.OvertimeRule{WeeklyHourLimit: 45}
	}

	summary := &EmployeeSummaryReport{
		EmployeeID: employeeID,
	}

	weeklyHours := make(map[int]float64)

	for _, e := range entries {
		hours := e.NetWorkingHours()
		summary.TotalHours += hours
		summary.WorkingDays++

		if summary.EmployeeName == "" {
			summary.EmployeeName = e.Employee.FirstName + " " + e.Employee.LastName
			summary.Department = e.Employee.Department.Name
		}

		_, week := e.ClockIn.ISOWeek()
		weeklyHours[week] += hours
	}

	// Calculate overtime from weekly totals
	for _, hours := range weeklyHours {
		if hours > rule.WeeklyHourLimit {
			summary.OvertimeHours += hours - rule.WeeklyHourLimit
		}
	}

	for _, l := range leaves {
		if l.Status == "approved" {
			summary.LeaveDays += l.TotalDays
		}
	}

	return summary, nil
}

// GetTrendAnalysis generates trend data across multiple months.
func (s *ReportingService) GetTrendAnalysis(deptID uint, startMonth, endMonth, year int) ([]TrendData, error) {
	var trends []TrendData

	for m := startMonth; m <= endMonth; m++ {
		start := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)

		var entries []core.TimeEntry
		var err error
		if deptID > 0 {
			entries, err = s.timeEntryRepo.ListByDepartment(deptID, start, end)
		} else {
			entries, err = s.timeEntryRepo.ListByDateRange(start, end)
		}
		if err != nil {
			return nil, err
		}

		rule, _ := s.overtimeRepo.GetActive()
		if rule == nil {
			rule = &core.OvertimeRule{WeeklyHourLimit: 45}
		}

		td := TrendData{Month: m, Year: year}

		weeklyHours := make(map[int]float64)
		for _, e := range entries {
			hours := e.NetWorkingHours()
			td.TotalHours += hours
			td.WorkingDays++

			_, week := e.ClockIn.ISOWeek()
			weeklyHours[week] += hours
		}

		for _, hours := range weeklyHours {
			if hours > rule.WeeklyHourLimit {
				td.OvertimeHours += hours - rule.WeeklyHourLimit
			}
		}

		// Count leave days
		var leaves []core.Leave
		if deptID > 0 {
			leaves, _ = s.leaveRepo.ListByDepartment(deptID, start, end)
		}
		for _, l := range leaves {
			if l.Status == "approved" {
				td.AbsenceDays += l.TotalDays
			}
		}

		trends = append(trends, td)
	}

	return trends, nil
}

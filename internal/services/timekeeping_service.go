package services

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

// TimekeepingService handles time entry operations including clock-in/out and overtime tracking.
type TimekeepingService struct {
	repo         repositories.TimeEntryRepositoryInterface
	scheduleRepo ScheduleRepositoryInterface
}

// NewTimekeepingService creates a new TimekeepingService.
func NewTimekeepingService(repo repositories.TimeEntryRepositoryInterface, scheduleRepo ScheduleRepositoryInterface) *TimekeepingService {
	return &TimekeepingService{repo: repo, scheduleRepo: scheduleRepo}
}

// ClockIn creates a new time entry with the current time as clock-in.
func (s *TimekeepingService) ClockIn(req core.ClockInRequest) (*core.TimeEntry, error) {
	// Check for already open entry
	existing, _ := s.repo.GetOpenEntry(req.EmployeeID)
	if existing != nil {
		return nil, fmt.Errorf("çalışanın zaten açık bir giriş kaydı var (ID: %d)", existing.ID)
	}

	now := time.Now()
	entry := &core.TimeEntry{
		EmployeeID: req.EmployeeID,
		ClockIn:    now,
		EntryType:  "normal",
		Source:     "auto",
		Notes:      req.Notes,
		Status:     "pending",
	}

	if err := s.repo.Create(entry); err != nil {
		return nil, fmt.Errorf("giriş kaydı oluşturulamadı: %w", err)
	}

	slog.Info("Clock-in recorded", "employeeID", req.EmployeeID, "clockIn", now)
	return entry, nil
}

// ClockOut finds the open time entry for the employee and sets the clock-out time.
func (s *TimekeepingService) ClockOut(req core.ClockOutRequest) (*core.TimeEntry, error) {
	entry, err := s.repo.GetOpenEntry(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("açık giriş kaydı bulunamadı: %w", err)
	}

	now := time.Now()
	if now.Before(entry.ClockIn) {
		return nil, fmt.Errorf("çıkış saati giriş saatinden önce olamaz")
	}

	entry.ClockOut = &now
	if req.Notes != "" {
		entry.Notes = req.Notes
	}

	// Classify entry type based on the day
	entry.EntryType = classifyEntryType(entry.ClockIn)

	if err := s.repo.Update(entry); err != nil {
		return nil, fmt.Errorf("çıkış kaydı güncellenemedi: %w", err)
	}

	slog.Info("Clock-out recorded",
		"employeeID", req.EmployeeID,
		"clockIn", entry.ClockIn,
		"clockOut", now,
		"netHours", entry.NetWorkingHours(),
	)
	return entry, nil
}

// CreateTimeEntry creates a manual time entry.
func (s *TimekeepingService) CreateTimeEntry(req core.TimeEntryRequest) (*core.TimeEntry, error) {
	if req.ClockOut != nil && req.ClockOut.Before(req.ClockIn) {
		return nil, fmt.Errorf("çıkış saati giriş saatinden önce olamaz")
	}

	entryType := req.EntryType
	if entryType == "" {
		entryType = classifyEntryType(req.ClockIn)
	}

	source := req.Source
	if source == "" {
		source = "manual"
	}

	entry := &core.TimeEntry{
		EmployeeID:   req.EmployeeID,
		ScheduleID:   req.ScheduleID,
		ClockIn:      req.ClockIn,
		ClockOut:     req.ClockOut,
		BreakMinutes: req.BreakMinutes,
		EntryType:    entryType,
		Source:       source,
		Notes:        req.Notes,
		Status:       "pending",
	}

	if err := s.repo.Create(entry); err != nil {
		return nil, err
	}

	slog.Info("Manual time entry created",
		"employeeID", req.EmployeeID,
		"clockIn", req.ClockIn,
		"source", source,
	)
	return entry, nil
}

// UpdateTimeEntry updates an existing time entry.
func (s *TimekeepingService) UpdateTimeEntry(id uint, req core.TimeEntryRequest) (*core.TimeEntry, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("kayıt bulunamadı: %w", err)
	}

	if req.ClockOut != nil && req.ClockOut.Before(req.ClockIn) {
		return nil, fmt.Errorf("çıkış saati giriş saatinden önce olamaz")
	}

	entry.EmployeeID = req.EmployeeID
	entry.ScheduleID = req.ScheduleID
	entry.ClockIn = req.ClockIn
	entry.ClockOut = req.ClockOut
	entry.BreakMinutes = req.BreakMinutes
	entry.Notes = req.Notes

	if req.EntryType != "" {
		entry.EntryType = req.EntryType
	}

	if err := s.repo.Update(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetPaginatedTimeEntries returns a paginated list of time entries with filters.
func (s *TimekeepingService) GetPaginatedTimeEntries(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) (*core.PaginationResult, error) {
	data, total, err := s.repo.ListPaginated(params, employeeID, departmentID, start, end)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}

	return &core.PaginationResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteTimeEntry removes a time entry.
func (s *TimekeepingService) DeleteTimeEntry(id uint) error {
	return s.repo.Delete(id)
}

// GetTimeEntry retrieves a single time entry by ID.
func (s *TimekeepingService) GetTimeEntry(id uint) (*core.TimeEntry, error) {
	return s.repo.GetByID(id)
}

// GetEmployeeTimeEntries returns time entries for an employee within a date range.
func (s *TimekeepingService) GetEmployeeTimeEntries(employeeID uint, start, end time.Time) ([]core.TimeEntry, error) {
	return s.repo.ListByEmployee(employeeID, start, end)
}

// GetDepartmentTimeEntries returns time entries for a department within a date range.
func (s *TimekeepingService) GetDepartmentTimeEntries(deptID uint, start, end time.Time) ([]core.TimeEntry, error) {
	return s.repo.ListByDepartment(deptID, start, end)
}

// GetTimeEntriesByDateRange returns all time entries within a date range.
func (s *TimekeepingService) GetTimeEntriesByDateRange(start, end time.Time) ([]core.TimeEntry, error) {
	return s.repo.ListByDateRange(start, end)
}

// GetPendingTimeEntries returns time entries with pending status.
func (s *TimekeepingService) GetPendingTimeEntries(start, end time.Time) ([]core.TimeEntry, error) {
	return s.repo.ListByStatus("pending", start, end)
}

// CalculateDailyHours returns the net working hours for a time entry.
func (s *TimekeepingService) CalculateDailyHours(entry *core.TimeEntry) float64 {
	return entry.NetWorkingHours()
}

// classifyEntryType determines the entry type based on the day of the week.
func classifyEntryType(t time.Time) string {
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return "weekend"
	}
	return "normal"
}

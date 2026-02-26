package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

// TimeEntryRepositoryInterface defines the contract for time entry data access.
type TimeEntryRepositoryInterface interface {
	Create(entry *core.TimeEntry) error
	Update(entry *core.TimeEntry) error
	GetByID(id uint) (*core.TimeEntry, error)
	Delete(id uint) error
	GetOpenEntry(employeeID uint) (*core.TimeEntry, error)
	ListByEmployee(employeeID uint, start, end time.Time) ([]core.TimeEntry, error)
	ListByDepartment(departmentID uint, start, end time.Time) ([]core.TimeEntry, error)
	ListByDateRange(start, end time.Time) ([]core.TimeEntry, error)
	ListByStatus(status string, start, end time.Time) ([]core.TimeEntry, error)
	ListPaginated(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) ([]core.TimeEntry, int64, error)
}

// TimeEntryRepository handles database operations for TimeEntry.
type TimeEntryRepository struct {
	db *gorm.DB
}

// NewTimeEntryRepository creates a new TimeEntryRepository.
func NewTimeEntryRepository(db *gorm.DB) *TimeEntryRepository {
	return &TimeEntryRepository{db: db}
}

func (r *TimeEntryRepository) Create(entry *core.TimeEntry) error {
	return r.db.Create(entry).Error
}

func (r *TimeEntryRepository) Update(entry *core.TimeEntry) error {
	return r.db.Save(entry).Error
}

func (r *TimeEntryRepository) GetByID(id uint) (*core.TimeEntry, error) {
	var entry core.TimeEntry
	err := r.db.Preload("Employee").Preload("Employee.Department").
		Preload("Schedule").Preload("Schedule.ShiftType").
		First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *TimeEntryRepository) Delete(id uint) error {
	return r.db.Delete(&core.TimeEntry{}, id).Error
}

// GetOpenEntry finds an open (no clock-out) entry for the given employee.
func (r *TimeEntryRepository) GetOpenEntry(employeeID uint) (*core.TimeEntry, error) {
	var entry core.TimeEntry
	err := r.db.Where("employee_id = ? AND clock_out IS NULL", employeeID).
		Order("clock_in DESC").First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *TimeEntryRepository) ListByEmployee(employeeID uint, start, end time.Time) ([]core.TimeEntry, error) {
	var entries []core.TimeEntry
	err := r.db.Preload("Employee").Preload("Schedule").Preload("Schedule.ShiftType").
		Where("employee_id = ? AND clock_in >= ? AND clock_in < ?", employeeID, start, end).
		Order("clock_in ASC").Find(&entries).Error
	return entries, err
}

func (r *TimeEntryRepository) ListByDepartment(departmentID uint, start, end time.Time) ([]core.TimeEntry, error) {
	var entries []core.TimeEntry
	err := r.db.Preload("Employee").Preload("Employee.Department").
		Preload("Schedule").Preload("Schedule.ShiftType").
		Joins("JOIN employees ON employees.id = time_entries.employee_id").
		Where("employees.department_id = ? AND time_entries.clock_in >= ? AND time_entries.clock_in < ?", departmentID, start, end).
		Order("time_entries.clock_in ASC").Find(&entries).Error
	return entries, err
}

func (r *TimeEntryRepository) ListByDateRange(start, end time.Time) ([]core.TimeEntry, error) {
	var entries []core.TimeEntry
	err := r.db.Preload("Employee").Preload("Employee.Department").
		Preload("Schedule").Preload("Schedule.ShiftType").
		Where("clock_in >= ? AND clock_in < ?", start, end).
		Order("clock_in ASC").Find(&entries).Error
	return entries, err
}

func (r *TimeEntryRepository) ListByStatus(status string, start, end time.Time) ([]core.TimeEntry, error) {
	var entries []core.TimeEntry
	err := r.db.Preload("Employee").Preload("Employee.Department").
		Where("status = ? AND clock_in >= ? AND clock_in < ?", status, start, end).
		Order("clock_in ASC").Find(&entries).Error
	return entries, err
}
func (r *TimeEntryRepository) ListPaginated(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) ([]core.TimeEntry, int64, error) {
	var entries []core.TimeEntry
	var total int64

	db := r.db.Model(&core.TimeEntry{}).Preload("Employee").Preload("Employee.Department").
		Preload("Schedule").Preload("Schedule.ShiftType")

	if employeeID > 0 {
		db = db.Where("employee_id = ?", employeeID)
	}
	if departmentID > 0 {
		db = db.Joins("JOIN employees ON employees.id = time_entries.employee_id").
			Where("employees.department_id = ?", departmentID)
	}
	if !start.IsZero() {
		db = db.Where("clock_in >= ?", start)
	}
	if !end.IsZero() {
		db = db.Where("clock_in < ?", end)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Joins("LEFT JOIN employees e ON e.id = time_entries.employee_id").
			Where("e.first_name ILIKE ? OR e.last_name ILIKE ?", search, search)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := db.Order("clock_in DESC").Offset(offset).Limit(params.Limit).Find(&entries).Error

	return entries, total, err
}

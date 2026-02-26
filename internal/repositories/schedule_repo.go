package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type ScheduleRepositoryInterface interface {
	Create(schedule *core.Schedule) error
	Update(schedule *core.Schedule) error
	GetByID(id uint) (*core.Schedule, error)
	Delete(id uint) error
	GetCombinedSchedule(month int, year int) ([]core.Schedule, error)
	DeleteByMonthYear(month int, year int) error
	ListPaginated(params core.PaginationParams, month, year int) ([]core.Schedule, int64, error)
}

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *core.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepository) Update(schedule *core.Schedule) error {
	return r.db.Save(schedule).Error
}

func (r *ScheduleRepository) GetByID(id uint) (*core.Schedule, error) {
	var schedule core.Schedule
	err := r.db.Preload("Employee").Preload("ShiftType").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) GetCombinedSchedule(month int, year int) ([]core.Schedule, error) {
	var schedules []core.Schedule
	// Filter by month and year based on Date field
	// Postgres generic way: extract(month from date) = ?
	// Or simpler: range query
	// But let's use a simple raw query or specific date range logic for better compatibility if GORM supports it
	// Actually, GORM specific dialect might be needed for extract.
	// Safer: start date <= date < end date
	// But month/year logic is passed as int.
	// Let's use SQl query.
	err := r.db.Preload("Employee").Preload("Employee.Department").Preload("Employee.Title").Preload("ShiftType").
		Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Order("date ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) DeleteByMonthYear(month int, year int) error {
	// Use Unscoped to permanently delete (not soft delete) old schedules
	return r.db.Unscoped().Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Delete(&core.Schedule{}).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&core.Schedule{}, id).Error
}

func (r *ScheduleRepository) ListPaginated(params core.PaginationParams, month, year int) ([]core.Schedule, int64, error) {
	var schedules []core.Schedule
	var total int64

	db := r.db.Model(&core.Schedule{}).Preload("Employee").Preload("ShiftType")

	if month > 0 {
		db = db.Where("EXTRACT(MONTH FROM date) = ?", month)
	}
	if year > 0 {
		db = db.Where("EXTRACT(YEAR FROM date) = ?", year)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Joins("LEFT JOIN employees ON employees.id = schedules.employee_id").
			Where("employees.first_name LIKE ? OR employees.last_name LIKE ?", search, search)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := db.Order("date DESC").Offset(offset).Limit(params.Limit).Find(&schedules).Error

	return schedules, total, err
}

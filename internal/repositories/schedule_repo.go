package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

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
	err := r.db.Preload("Employee").Preload("ShiftType").
		Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Order("date ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) DeleteByMonthYear(month int, year int) error {
	return r.db.Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Delete(&core.Schedule{}).Error
}

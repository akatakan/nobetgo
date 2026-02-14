package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type AttendanceRepositoryInterface interface {
	Create(attendance *core.Attendance) error
	GetCombinedReport(month int, year int) ([]core.Attendance, error)
}

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *core.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) GetCombinedReport(month int, year int) ([]core.Attendance, error) {
	var report []core.Attendance
	// Join with Schedule to filter by date
	err := r.db.Preload("Schedule").Preload("Schedule.Employee").Preload("Schedule.ShiftType").
		Joins("JOIN schedules ON schedules.id = attendances.schedule_id").
		Where("EXTRACT(MONTH FROM schedules.date) = ? AND EXTRACT(YEAR FROM schedules.date) = ?", month, year).
		Find(&report).Error
	return report, err
}

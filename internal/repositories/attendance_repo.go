package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type AttendanceRepositoryInterface interface {
	Create(attendance *core.Attendance) error
	Update(attendance *core.Attendance) error
	GetByID(id uint) (*core.Attendance, error)
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

func (r *AttendanceRepository) Update(attendance *core.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *AttendanceRepository) GetByID(id uint) (*core.Attendance, error) {
	var att core.Attendance
	err := r.db.Preload("Schedule").Preload("Schedule.ShiftType").First(&att, id).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
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

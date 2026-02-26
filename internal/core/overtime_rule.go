package core

import (
	"time"

	"gorm.io/gorm"
)

// OvertimeRule defines overtime calculation rules (multipliers for weekday, weekend, holiday, night).
type OvertimeRule struct {
	gorm.Model
	Name               string  `gorm:"not null;uniqueIndex" json:"name"`
	WeeklyHourLimit    float64 `gorm:"default:45" json:"weekly_hour_limit"`
	DailyHourLimit     float64 `gorm:"default:11" json:"daily_hour_limit"`
	OvertimeMultiplier float64 `gorm:"default:1.5" json:"overtime_multiplier"`
	WeekendMultiplier  float64 `gorm:"default:2.0" json:"weekend_multiplier"`
	HolidayMultiplier  float64 `gorm:"default:2.5" json:"holiday_multiplier"`
	NightShiftExtra    float64 `gorm:"default:0.1" json:"night_shift_extra"`
	IsActive           bool    `gorm:"default:true" json:"is_active"`
}

// PublicHoliday represents an official public holiday date.
type PublicHoliday struct {
	gorm.Model
	Name string    `gorm:"not null" json:"name"`
	Date time.Time `gorm:"not null;uniqueIndex" json:"date"`
}

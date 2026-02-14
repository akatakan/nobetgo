package core

import (
	"time"

	"gorm.io/gorm"
)

// Attendance represents the actual realized shift
type Attendance struct {
	gorm.Model
	ScheduleID      uint      `gorm:"not null"` // Links to the planned schedule
	Schedule        Schedule  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ActualStartTime time.Time `gorm:"not null"`
	ActualEndTime   time.Time `gorm:"not null"`
	Notes           string
	IsOvertime      bool    `gorm:"default:false"`
	OvertimeHours   float64 `gorm:"default:0"`
}

// AttendanceRequest for logging time
type AttendanceRequest struct {
	ScheduleID      uint      `json:"schedule_id" binding:"required"`
	ActualStartTime time.Time `json:"actual_start_time" binding:"required"`
	ActualEndTime   time.Time `json:"actual_end_time" binding:"required"`
	Notes           string    `json:"notes"`
}

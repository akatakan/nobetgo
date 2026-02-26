package core

import (
	"time"

	"gorm.io/gorm"
)

// TimeEntry represents a single clock-in/clock-out record for an employee.
type TimeEntry struct {
	gorm.Model
	EmployeeID   uint       `gorm:"not null;index" json:"employee_id"`
	Employee     Employee   `json:"employee,omitempty"`
	ScheduleID   *uint      `gorm:"index" json:"schedule_id,omitempty"`
	Schedule     *Schedule  `json:"schedule,omitempty"`
	ClockIn      time.Time  `gorm:"not null" json:"clock_in"`
	ClockOut     *time.Time `json:"clock_out,omitempty"`
	BreakMinutes int        `gorm:"default:0" json:"break_minutes"`
	EntryType    string     `gorm:"not null;default:'normal'" json:"entry_type"` // normal, overtime, holiday, weekend
	Source       string     `gorm:"default:'manual'" json:"source"`             // manual, auto, import
	Notes        string     `json:"notes"`
	Status       string     `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	ApprovedBy   *uint      `json:"approved_by,omitempty"`
}

// TimeEntryRequest is used for creating or updating a time entry.
type TimeEntryRequest struct {
	EmployeeID   uint       `json:"employee_id" binding:"required"`
	ScheduleID   *uint      `json:"schedule_id,omitempty"`
	ClockIn      time.Time  `json:"clock_in" binding:"required"`
	ClockOut     *time.Time `json:"clock_out,omitempty"`
	BreakMinutes int        `json:"break_minutes"`
	EntryType    string     `json:"entry_type"`
	Source       string     `json:"source"`
	Notes        string     `json:"notes"`
}

// ClockInRequest is used for automatic clock-in.
type ClockInRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Notes      string `json:"notes"`
}

// ClockOutRequest is used for automatic clock-out.
type ClockOutRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Notes      string `json:"notes"`
}

// NetWorkingHours returns the net working hours for a completed time entry.
func (t *TimeEntry) NetWorkingHours() float64 {
	if t.ClockOut == nil {
		return 0
	}
	total := t.ClockOut.Sub(t.ClockIn).Hours()
	breakHours := float64(t.BreakMinutes) / 60.0
	net := total - breakHours
	if net < 0 {
		return 0
	}
	return net
}

package core

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Date        time.Time `gorm:"not null"`
	EmployeeID  uint      `gorm:"not null"`
	ShiftTypeID uint      `gorm:"not null"`
	Employee    Employee
	ShiftType   ShiftType
	IsLocked    bool `gorm:"default:false"` // If true, this assignment cannot be changed by auto-scheduler
}

// ScheduleRequest represents parameters for generating a schedule
type ScheduleRequest struct {
	Month              int     `json:"month" binding:"required,min=1,max=12"`
	Year               int     `json:"year" binding:"required,min=2024"`
	DepartmentID       uint    `json:"department_id"`
	ShiftTypeIDs       []uint  `json:"shift_type_ids"`
	EmployeeIDs        []uint  `json:"employee_ids"`
	OvertimeThreshold  float64 `json:"overtime_threshold"`  // Default 45.0
	OvertimeMultiplier float64 `json:"overtime_multiplier"` // Default 1.5
}

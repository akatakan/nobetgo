package core

import (
	"time"

	"gorm.io/gorm"
)

// LeaveType defines categories of leave (annual, sick, administrative, etc.).
type LeaveType struct {
	gorm.Model
	Name             string `gorm:"not null;uniqueIndex" json:"name"`
	DefaultDays      int    `gorm:"default:0" json:"default_days"`
	IsPaid           bool   `gorm:"default:true" json:"is_paid"`
	RequiresApproval bool   `gorm:"default:true" json:"requires_approval"`
	Color            string `json:"color"`
}

// Leave represents a leave or absence record for an employee.
type Leave struct {
	gorm.Model
	EmployeeID  uint       `gorm:"not null;index" json:"employee_id"`
	Employee    Employee   `json:"employee,omitempty"`
	LeaveTypeID uint       `gorm:"not null" json:"leave_type_id"`
	LeaveType   LeaveType  `json:"leave_type,omitempty"`
	StartDate   time.Time  `gorm:"not null" json:"start_date"`
	EndDate     time.Time  `gorm:"not null" json:"end_date"`
	TotalDays   float64    `gorm:"not null" json:"total_days"`
	Reason      string     `json:"reason"`
	Status      string     `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	ApprovedBy  *uint      `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
}

// LeaveBalance tracks available leave days per employee, per leave type, per year.
type LeaveBalance struct {
	gorm.Model
	EmployeeID    uint      `gorm:"not null;index;uniqueIndex:idx_leave_balance" json:"employee_id"`
	Employee      Employee  `json:"employee,omitempty"`
	LeaveTypeID   uint      `gorm:"not null;uniqueIndex:idx_leave_balance" json:"leave_type_id"`
	LeaveType     LeaveType `json:"leave_type,omitempty"`
	Year          int       `gorm:"not null;uniqueIndex:idx_leave_balance" json:"year"`
	TotalDays     float64   `gorm:"default:0" json:"total_days"`
	UsedDays      float64   `gorm:"default:0" json:"used_days"`
	RemainingDays float64   `gorm:"default:0" json:"remaining_days"`
}

// LeaveRequest is used for creating a new leave request.
type LeaveRequest struct {
	EmployeeID  uint      `json:"employee_id" binding:"required"`
	LeaveTypeID uint      `json:"leave_type_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Reason      string    `json:"reason"`
}

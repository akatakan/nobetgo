package core

import "gorm.io/gorm"

// ShiftType defines a work shift pattern with optional rotation and department binding.
type ShiftType struct {
	gorm.Model
	Name         string      `gorm:"not null;uniqueIndex" json:"Name"`
	Description  string      `json:"Description"`
	StartTime    string      `json:"StartTime"` // HH:mm format
	EndTime      string      `json:"EndTime"`   // HH:mm format
	Color        string      `json:"Color"`     // Hex color code for UI
	BreakMinutes int         `gorm:"default:0" json:"BreakMinutes"`
	IsNightShift bool        `gorm:"default:false" json:"IsNightShift"`
	RotationDays int         `gorm:"default:0" json:"RotationDays"` // 0 = no rotation
	DepartmentID *uint       `gorm:"index" json:"DepartmentID,omitempty"`
	Department   *Department `json:"Department,omitempty"`
}

// RotationPlan defines a shift rotation schedule for a department.
type RotationPlan struct {
	gorm.Model
	Name         string     `gorm:"not null" json:"name"`
	DepartmentID uint       `gorm:"not null;index" json:"department_id"`
	Department   Department `json:"department,omitempty"`
	ShiftTypeID  uint       `gorm:"not null" json:"shift_type_id"`
	ShiftType    ShiftType  `json:"shift_type,omitempty"`
	CycleDays    int        `gorm:"not null" json:"cycle_days"`
	RestDays     int        `gorm:"default:1" json:"rest_days"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
}

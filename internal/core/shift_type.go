package core

import "gorm.io/gorm"

type ShiftType struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex"`
	Description string
	StartTime   string // HH:mm format
	EndTime     string // HH:mm format
	Color       string // Hex color code for UI
}

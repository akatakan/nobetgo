package core

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"Name"`
	Floor       int    `gorm:"not null;default:1" json:"Floor"`
	Description string `json:"Description"`
	BedCapacity int    `gorm:"default:0" json:"BedCapacity"` // Added for bed-capacity based scheduling
}

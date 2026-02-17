package core

import "gorm.io/gorm"

type Title struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex" json:"Name"`
}

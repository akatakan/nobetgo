package core

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Title      string
	Department string
	Email      string `gorm:"uniqueIndex"`
	Phone      string
	HourlyRate float64 `gorm:"default:0"`
	IsActive   bool    `gorm:"default:true"`
}

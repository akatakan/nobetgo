package core

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	FirstName     string     `gorm:"not null" json:"FirstName"`
	LastName      string     `gorm:"not null" json:"LastName"`
	TitleID       uint       `json:"TitleID"`
	Title         Title      `json:"Title"`
	DepartmentID  uint       `json:"DepartmentID"`
	Department    Department `json:"Department"`
	Email         string     `json:"Email"`
	Phone         string     `json:"Phone"`
	HourlyRate    float64    `gorm:"default:0" json:"HourlyRate"`
	IsShiftWorker bool       `gorm:"default:true" json:"IsShiftWorker"` // If false, excluded from auto-schedule
	IsActive      bool       `gorm:"default:true" json:"IsActive"`
}

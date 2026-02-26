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
	Username      string     `gorm:"type:varchar(50);uniqueIndex" json:"Username"`
	Email         string     `json:"Email"`
	Phone         string     `json:"Phone"`
	HourlyRate    float64    `gorm:"default:0" json:"HourlyRate"`
	IsShiftWorker bool       `gorm:"default:true" json:"IsShiftWorker"` // If false, excluded from auto-schedule
	IsActive      bool       `gorm:"default:true" json:"IsActive"`
	Competencies  string     `gorm:"type:text" json:"Competencies"` // Stored as comma-separated string or JSON string to keep it DB agnostic (SQLite/Postgres)
	FatigueScore  int        `gorm:"default:0" json:"FatigueScore"`
	HeroPoint     int        `gorm:"default:0" json:"HeroPoint"`
	PasswordHash  string     `gorm:"type:varchar(255)" json:"-"`
	Role          string     `gorm:"type:varchar(20);default:'user'" json:"Role"` // admin or user
}

package core

import "gorm.io/gorm"

// Notification represents an in-app system notification for users
type Notification struct {
	gorm.Model
	EmployeeID  uint     `gorm:"not null;index" json:"employee_id"`
	Employee    Employee `json:"employee,omitempty"`
	Title       string   `gorm:"not null" json:"title"`
	Message     string   `gorm:"not null" json:"message"`
	Type        string   `gorm:"not null" json:"type"` // Example: info, warning, shift_change, approval
	IsRead      bool     `gorm:"default:false;index" json:"is_read"`
	ActionURL   string   `json:"action_url,omitempty"`   // Optional link to navigate to when clicked
	RelatedType string   `json:"related_type,omitempty"` // Entity type (e.g. "Leave", "Schedule")
	RelatedID   uint     `json:"related_id,omitempty"`   // Entity ID
}

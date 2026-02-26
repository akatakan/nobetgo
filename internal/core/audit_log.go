package core

import "gorm.io/gorm"

// AuditLog records all changes to tracked entities for full traceability.
type AuditLog struct {
	gorm.Model
	EntityType  string `gorm:"not null;index" json:"entity_type"` // time_entry, leave, schedule
	EntityID    uint   `gorm:"not null;index" json:"entity_id"`
	Action      string `gorm:"not null" json:"action"` // create, update, delete, approve, reject
	FieldName   string `json:"field_name,omitempty"`
	OldValue    string `json:"old_value,omitempty"`
	NewValue    string `json:"new_value,omitempty"`
	PerformedBy uint   `gorm:"not null" json:"performed_by"`
	IPAddress   string `json:"ip_address,omitempty"`
}

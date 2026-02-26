package core

import (
	"time"

	"gorm.io/gorm"
)

// PasswordResetToken represents a token used for password reset requests.
type PasswordResetToken struct {
	gorm.Model
	EmployeeID uint       `json:"employee_id" gorm:"index"`
	Employee   Employee   `json:"-"`
	Token      string     `json:"token" gorm:"uniqueIndex"`
	ExpiresAt  time.Time  `json:"expires_at"`
	UsedAt     *time.Time `json:"used_at"`
}

func (t *PasswordResetToken) IsValid() bool {
	return t.UsedAt == nil && t.ExpiresAt.After(time.Now())
}

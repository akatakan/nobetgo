package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type PasswordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

func (r *PasswordResetTokenRepository) Create(token *core.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *PasswordResetTokenRepository) GetByToken(token string) (*core.PasswordResetToken, error) {
	var t core.PasswordResetToken
	err := r.db.Preload("Employee").Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PasswordResetTokenRepository) MarkAsUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&core.PasswordResetToken{}).Where("id = ?", id).Update("used_at", &now).Error
}

func (r *PasswordResetTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&core.PasswordResetToken{}).Error
}

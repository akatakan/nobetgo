package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

// AuditLogRepositoryInterface defines the contract for audit log data access.
type AuditLogRepositoryInterface interface {
	Create(log *core.AuditLog) error
	ListByEntity(entityType string, entityID uint) ([]core.AuditLog, error)
	ListByDateRange(start, end time.Time) ([]core.AuditLog, error)
	ListByPerformer(performerID uint, start, end time.Time) ([]core.AuditLog, error)
}

// AuditLogRepository handles database operations for AuditLog.
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new AuditLogRepository.
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *core.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) ListByEntity(entityType string, entityID uint) ([]core.AuditLog, error) {
	var logs []core.AuditLog
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) ListByDateRange(start, end time.Time) ([]core.AuditLog, error) {
	var logs []core.AuditLog
	err := r.db.Where("created_at >= ? AND created_at < ?", start, end).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) ListByPerformer(performerID uint, start, end time.Time) ([]core.AuditLog, error) {
	var logs []core.AuditLog
	err := r.db.Where("performed_by = ? AND created_at >= ? AND created_at < ?", performerID, start, end).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

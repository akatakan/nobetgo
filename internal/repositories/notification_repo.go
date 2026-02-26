package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(notification *core.Notification) error
	GetUnreadByEmployee(employeeID uint) ([]core.Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(employeeID uint) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *core.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetUnreadByEmployee(employeeID uint) ([]core.Notification, error) {
	var notifications []core.Notification
	err := r.db.Where("employee_id = ? AND is_read = ?", employeeID, false).
		Order("created_at desc").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&core.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(employeeID uint) error {
	return r.db.Model(&core.Notification{}).Where("employee_id = ? AND is_read = ?", employeeID, false).Update("is_read", true).Error
}

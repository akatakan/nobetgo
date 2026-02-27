package services

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

type NotificationService struct {
	repo repositories.NotificationRepositoryInterface
}

func NewNotificationService(repo repositories.NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(employeeID uint, title, message, notifType, url, relType string, relID uint) error {
	notif := &core.Notification{
		EmployeeID:  employeeID,
		Title:       title,
		Message:     message,
		Type:        notifType,
		ActionURL:   url,
		RelatedType: relType,
		RelatedID:   relID,
	}
	return s.repo.Create(notif)
}

func (s *NotificationService) GetUnread(employeeID uint) ([]core.Notification, error) {
	return s.repo.GetUnreadByEmployee(employeeID)
}

func (s *NotificationService) MarkAsRead(id uint, employeeID uint) (bool, error) {
	return s.repo.MarkAsReadForEmployee(id, employeeID)
}

func (s *NotificationService) MarkAllAsRead(employeeID uint) error {
	return s.repo.MarkAllAsRead(employeeID)
}

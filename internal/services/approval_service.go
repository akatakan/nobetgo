package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

// ApprovalService handles multi-step approval workflows and audit logging.
type ApprovalService struct {
	auditRepo     repositories.AuditLogRepositoryInterface
	timeEntryRepo repositories.TimeEntryRepositoryInterface
	leaveRepo     repositories.LeaveRepositoryInterface
}

// NewApprovalService creates a new ApprovalService.
func NewApprovalService(
	auditRepo repositories.AuditLogRepositoryInterface,
	timeEntryRepo repositories.TimeEntryRepositoryInterface,
	leaveRepo repositories.LeaveRepositoryInterface,
) *ApprovalService {
	return &ApprovalService{
		auditRepo:     auditRepo,
		timeEntryRepo: timeEntryRepo,
		leaveRepo:     leaveRepo,
	}
}

// ApproveTimeEntry approves a pending time entry.
func (s *ApprovalService) ApproveTimeEntry(id uint, approverID uint) (*core.TimeEntry, error) {
	entry, err := s.timeEntryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("kayıt bulunamadı: %w", err)
	}

	if entry.Status != "pending" {
		return nil, fmt.Errorf("sadece bekleyen kayıtlar onaylanabilir (mevcut durum: %s)", entry.Status)
	}

	oldStatus := entry.Status
	entry.Status = "approved"
	entry.ApprovedBy = &approverID

	if err := s.timeEntryRepo.Update(entry); err != nil {
		return nil, err
	}

	s.logChange("time_entry", id, "approve", "status", oldStatus, "approved", approverID)
	slog.Info("Time entry approved", "entryID", id, "approverID", approverID)
	return entry, nil
}

// RejectTimeEntry rejects a pending time entry.
func (s *ApprovalService) RejectTimeEntry(id uint, approverID uint) (*core.TimeEntry, error) {
	entry, err := s.timeEntryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("kayıt bulunamadı: %w", err)
	}

	if entry.Status != "pending" {
		return nil, fmt.Errorf("sadece bekleyen kayıtlar reddedilebilir (mevcut durum: %s)", entry.Status)
	}

	oldStatus := entry.Status
	entry.Status = "rejected"
	entry.ApprovedBy = &approverID

	if err := s.timeEntryRepo.Update(entry); err != nil {
		return nil, err
	}

	s.logChange("time_entry", id, "reject", "status", oldStatus, "rejected", approverID)
	slog.Info("Time entry rejected", "entryID", id, "approverID", approverID)
	return entry, nil
}

// GetPendingApprovals returns all pending time entries and leaves.
type PendingApprovals struct {
	TimeEntries []core.TimeEntry `json:"time_entries"`
	Leaves      []core.Leave     `json:"leaves"`
}

// GetPendingApprovals retrieves all pending items for review.
func (s *ApprovalService) GetPendingApprovals() (*PendingApprovals, error) {
	// Get pending time entries (last 90 days)
	now := time.Now()
	start := now.AddDate(0, -3, 0)
	entries, err := s.timeEntryRepo.ListByStatus("pending", start, now)
	if err != nil {
		return nil, err
	}

	// Get pending leaves
	leaves, err := s.leaveRepo.ListByStatus("pending")
	if err != nil {
		return nil, err
	}

	return &PendingApprovals{
		TimeEntries: entries,
		Leaves:      leaves,
	}, nil
}

// GetAuditLogs returns audit logs for a specific entity.
func (s *ApprovalService) GetAuditLogs(entityType string, entityID uint) ([]core.AuditLog, error) {
	return s.auditRepo.ListByEntity(entityType, entityID)
}

// GetAuditLogsByDateRange returns audit logs within a date range.
func (s *ApprovalService) GetAuditLogsByDateRange(start, end time.Time) ([]core.AuditLog, error) {
	return s.auditRepo.ListByDateRange(start, end)
}

// LogChange creates an audit log entry.
func (s *ApprovalService) LogChange(entityType string, entityID uint, action, field, oldVal, newVal string, performedBy uint) {
	s.logChange(entityType, entityID, action, field, oldVal, newVal, performedBy)
}

func (s *ApprovalService) logChange(entityType string, entityID uint, action, field, oldVal, newVal string, performedBy uint) {
	log := &core.AuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		FieldName:   field,
		OldValue:    oldVal,
		NewValue:    newVal,
		PerformedBy: performedBy,
	}

	if err := s.auditRepo.Create(log); err != nil {
		slog.Error("Failed to create audit log",
			"entityType", entityType,
			"entityID", entityID,
			"action", action,
			"error", err,
		)
	}
}

package services

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

// LeaveService handles leave request, approval, and balance management.
type LeaveService struct {
	repo repositories.LeaveRepositoryInterface
}

// NewLeaveService creates a new LeaveService.
func NewLeaveService(repo repositories.LeaveRepositoryInterface) *LeaveService {
	return &LeaveService{repo: repo}
}

// RequestLeave creates a new leave request after validation.
func (s *LeaveService) RequestLeave(req core.LeaveRequest) (*core.Leave, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("bitiş tarihi başlangıç tarihinden önce olamaz")
	}

	// Check for overlapping leaves
	overlap, err := s.repo.HasOverlap(req.EmployeeID, req.StartDate, req.EndDate, 0)
	if err != nil {
		return nil, fmt.Errorf("çakışma kontrolü yapılamadı: %w", err)
	}
	if overlap {
		return nil, fmt.Errorf("bu tarih aralığında zaten bir izin kaydı mevcut")
	}

	totalDays := calculateBusinessDays(req.StartDate, req.EndDate)

	leave := &core.Leave{
		EmployeeID:  req.EmployeeID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalDays:   totalDays,
		Reason:      req.Reason,
		Status:      "pending",
	}

	if err := s.repo.Create(leave); err != nil {
		return nil, err
	}

	slog.Info("Leave requested",
		"employeeID", req.EmployeeID,
		"startDate", req.StartDate,
		"endDate", req.EndDate,
		"totalDays", totalDays,
	)
	return leave, nil
}

// ApproveLeave approves a pending leave request and updates the balance.
func (s *LeaveService) ApproveLeave(id uint, approverID uint) (*core.Leave, error) {
	leave, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("izin kaydı bulunamadı: %w", err)
	}

	if leave.Status != "pending" {
		return nil, fmt.Errorf("sadece bekleyen izinler onaylanabilir (mevcut durum: %s)", leave.Status)
	}

	now := time.Now()
	leave.Status = "approved"
	leave.ApprovedBy = &approverID
	leave.ApprovedAt = &now

	if err := s.repo.Update(leave); err != nil {
		return nil, err
	}

	// Update balance
	year := leave.StartDate.Year()
	balance, err := s.repo.GetBalance(leave.EmployeeID, leave.LeaveTypeID, year)
	if err == nil && balance != nil {
		balance.UsedDays += leave.TotalDays
		balance.RemainingDays = balance.TotalDays - balance.UsedDays
		_ = s.repo.UpsertBalance(balance)
	}

	slog.Info("Leave approved", "leaveID", id, "approverID", approverID)
	return leave, nil
}

// RejectLeave rejects a pending leave request.
func (s *LeaveService) RejectLeave(id uint, approverID uint) (*core.Leave, error) {
	leave, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("izin kaydı bulunamadı: %w", err)
	}

	if leave.Status != "pending" {
		return nil, fmt.Errorf("sadece bekleyen izinler reddedilebilir (mevcut durum: %s)", leave.Status)
	}

	now := time.Now()
	leave.Status = "rejected"
	leave.ApprovedBy = &approverID
	leave.ApprovedAt = &now

	if err := s.repo.Update(leave); err != nil {
		return nil, err
	}

	slog.Info("Leave rejected", "leaveID", id, "approverID", approverID)
	return leave, nil
}

// GetLeave retrieves a single leave by ID.
func (s *LeaveService) GetLeave(id uint) (*core.Leave, error) {
	return s.repo.GetByID(id)
}

// GetEmployeeLeaves returns all leaves for an employee within a date range.
func (s *LeaveService) GetEmployeeLeaves(employeeID uint, start, end time.Time) ([]core.Leave, error) {
	return s.repo.ListByEmployee(employeeID, start, end)
}

// GetDepartmentLeaves returns all leaves for a department within a date range.
func (s *LeaveService) GetDepartmentLeaves(deptID uint, start, end time.Time) ([]core.Leave, error) {
	return s.repo.ListByDepartment(deptID, start, end)
}

// GetPendingLeaves returns all pending leave requests.
func (s *LeaveService) GetPendingLeaves() ([]core.Leave, error) {
	return s.repo.ListByStatus("pending")
}

// GetLeaveBalance returns balances for an employee in a given year.
func (s *LeaveService) GetLeaveBalance(employeeID uint, year int) ([]core.LeaveBalance, error) {
	return s.repo.GetAllBalances(employeeID, year)
}

// GetPaginatedLeaves returns a paginated list of leaves with filters.
func (s *LeaveService) GetPaginatedLeaves(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) (*core.PaginationResult, error) {
	data, total, err := s.repo.ListPaginated(params, employeeID, departmentID, start, end)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &core.PaginationResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// InitializeBalance creates initial leave balances for an employee based on leave type defaults.
func (s *LeaveService) InitializeBalance(employeeID uint, year int) error {
	leaveTypes, err := s.repo.ListLeaveTypes()
	if err != nil {
		return err
	}

	for _, lt := range leaveTypes {
		balance := &core.LeaveBalance{
			EmployeeID:    employeeID,
			LeaveTypeID:   lt.ID,
			Year:          year,
			TotalDays:     float64(lt.DefaultDays),
			UsedDays:      0,
			RemainingDays: float64(lt.DefaultDays),
		}
		if err := s.repo.UpsertBalance(balance); err != nil {
			slog.Warn("Failed to initialize balance",
				"employeeID", employeeID,
				"leaveTypeID", lt.ID,
				"error", err,
			)
		}
	}

	slog.Info("Leave balances initialized", "employeeID", employeeID, "year", year)
	return nil
}

// --- LeaveType CRUD ---

// CreateLeaveType creates a new leave type.
func (s *LeaveService) CreateLeaveType(lt *core.LeaveType) error {
	return s.repo.CreateLeaveType(lt)
}

// UpdateLeaveType updates an existing leave type.
func (s *LeaveService) UpdateLeaveType(lt *core.LeaveType) error {
	return s.repo.UpdateLeaveType(lt)
}

// GetLeaveType retrieves a leave type by ID.
func (s *LeaveService) GetLeaveType(id uint) (*core.LeaveType, error) {
	return s.repo.GetLeaveTypeByID(id)
}

// GetAllLeaveTypes returns all leave types.
func (s *LeaveService) GetAllLeaveTypes() ([]core.LeaveType, error) {
	return s.repo.ListLeaveTypes()
}

// DeleteLeaveType removes a leave type.
func (s *LeaveService) DeleteLeaveType(id uint) error {
	return s.repo.DeleteLeaveType(id)
}

// calculateBusinessDays counts working days between two dates (excluding weekends).
func calculateBusinessDays(start, end time.Time) float64 {
	days := 0.0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			days++
		}
	}
	return math.Max(days, 0)
}

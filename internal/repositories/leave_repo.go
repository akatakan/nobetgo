package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

// LeaveRepositoryInterface defines the contract for leave data access.
type LeaveRepositoryInterface interface {
	Create(leave *core.Leave) error
	Update(leave *core.Leave) error
	GetByID(id uint) (*core.Leave, error)
	Delete(id uint) error
	ListByEmployee(employeeID uint, start, end time.Time) ([]core.Leave, error)
	ListByDepartment(departmentID uint, start, end time.Time) ([]core.Leave, error)
	ListByStatus(status string) ([]core.Leave, error)
	HasOverlap(employeeID uint, start, end time.Time, excludeID uint) (bool, error)

	// LeaveType CRUD
	CreateLeaveType(lt *core.LeaveType) error
	UpdateLeaveType(lt *core.LeaveType) error
	GetLeaveTypeByID(id uint) (*core.LeaveType, error)
	ListLeaveTypes() ([]core.LeaveType, error)
	DeleteLeaveType(id uint) error

	// LeaveBalance
	GetBalance(employeeID uint, leaveTypeID uint, year int) (*core.LeaveBalance, error)
	GetAllBalances(employeeID uint, year int) ([]core.LeaveBalance, error)
	UpsertBalance(balance *core.LeaveBalance) error
}

// LeaveRepository handles database operations for Leave, LeaveType, and LeaveBalance.
type LeaveRepository struct {
	db *gorm.DB
}

// NewLeaveRepository creates a new LeaveRepository.
func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// --- Leave CRUD ---

func (r *LeaveRepository) Create(leave *core.Leave) error {
	return r.db.Create(leave).Error
}

func (r *LeaveRepository) Update(leave *core.Leave) error {
	return r.db.Save(leave).Error
}

func (r *LeaveRepository) GetByID(id uint) (*core.Leave, error) {
	var leave core.Leave
	err := r.db.Preload("Employee").Preload("LeaveType").First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *LeaveRepository) Delete(id uint) error {
	return r.db.Delete(&core.Leave{}, id).Error
}

func (r *LeaveRepository) ListByEmployee(employeeID uint, start, end time.Time) ([]core.Leave, error) {
	var leaves []core.Leave
	err := r.db.Preload("LeaveType").
		Where("employee_id = ? AND start_date >= ? AND start_date < ?", employeeID, start, end).
		Order("start_date ASC").Find(&leaves).Error
	return leaves, err
}

func (r *LeaveRepository) ListByDepartment(departmentID uint, start, end time.Time) ([]core.Leave, error) {
	var leaves []core.Leave
	err := r.db.Preload("Employee").Preload("LeaveType").
		Joins("JOIN employees ON employees.id = leaves.employee_id").
		Where("employees.department_id = ? AND leaves.start_date >= ? AND leaves.start_date < ?", departmentID, start, end).
		Order("leaves.start_date ASC").Find(&leaves).Error
	return leaves, err
}

func (r *LeaveRepository) ListByStatus(status string) ([]core.Leave, error) {
	var leaves []core.Leave
	err := r.db.Preload("Employee").Preload("LeaveType").
		Where("status = ?", status).
		Order("start_date ASC").Find(&leaves).Error
	return leaves, err
}

// HasOverlap checks if the employee already has an approved/pending leave overlapping the given range.
func (r *LeaveRepository) HasOverlap(employeeID uint, start, end time.Time, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&core.Leave{}).
		Where("employee_id = ? AND status != 'rejected' AND start_date < ? AND end_date > ?", employeeID, end, start)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// --- LeaveType CRUD ---

func (r *LeaveRepository) CreateLeaveType(lt *core.LeaveType) error {
	return r.db.Create(lt).Error
}

func (r *LeaveRepository) UpdateLeaveType(lt *core.LeaveType) error {
	return r.db.Save(lt).Error
}

func (r *LeaveRepository) GetLeaveTypeByID(id uint) (*core.LeaveType, error) {
	var lt core.LeaveType
	if err := r.db.First(&lt, id).Error; err != nil {
		return nil, err
	}
	return &lt, nil
}

func (r *LeaveRepository) ListLeaveTypes() ([]core.LeaveType, error) {
	var types []core.LeaveType
	err := r.db.Order("name ASC").Find(&types).Error
	return types, err
}

func (r *LeaveRepository) DeleteLeaveType(id uint) error {
	return r.db.Delete(&core.LeaveType{}, id).Error
}

// --- LeaveBalance ---

func (r *LeaveRepository) GetBalance(employeeID uint, leaveTypeID uint, year int) (*core.LeaveBalance, error) {
	var balance core.LeaveBalance
	err := r.db.Where("employee_id = ? AND leave_type_id = ? AND year = ?", employeeID, leaveTypeID, year).
		First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *LeaveRepository) GetAllBalances(employeeID uint, year int) ([]core.LeaveBalance, error) {
	var balances []core.LeaveBalance
	err := r.db.Preload("LeaveType").
		Where("employee_id = ? AND year = ?", employeeID, year).
		Find(&balances).Error
	return balances, err
}

// UpsertBalance creates or updates a leave balance record.
func (r *LeaveRepository) UpsertBalance(balance *core.LeaveBalance) error {
	var existing core.LeaveBalance
	err := r.db.Where("employee_id = ? AND leave_type_id = ? AND year = ?",
		balance.EmployeeID, balance.LeaveTypeID, balance.Year).First(&existing).Error
	if err != nil {
		// Not found — create
		return r.db.Create(balance).Error
	}
	existing.TotalDays = balance.TotalDays
	existing.UsedDays = balance.UsedDays
	existing.RemainingDays = balance.RemainingDays
	return r.db.Save(&existing).Error
}

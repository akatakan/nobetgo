package repositories

import (
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

// OvertimeRuleRepositoryInterface defines the contract for overtime rule data access.
type OvertimeRuleRepositoryInterface interface {
	Create(rule *core.OvertimeRule) error
	Update(rule *core.OvertimeRule) error
	GetByID(id uint) (*core.OvertimeRule, error)
	GetActive() (*core.OvertimeRule, error)
	List() ([]core.OvertimeRule, error)
	Delete(id uint) error

	// PublicHoliday
	CreateHoliday(holiday *core.PublicHoliday) error
	UpdateHoliday(holiday *core.PublicHoliday) error
	GetHolidayByID(id uint) (*core.PublicHoliday, error)
	ListHolidays(year int) ([]core.PublicHoliday, error)
	DeleteHoliday(id uint) error
	IsHoliday(date time.Time) (bool, error)
}

// OvertimeRuleRepository handles database operations for OvertimeRule and PublicHoliday.
type OvertimeRuleRepository struct {
	db *gorm.DB
}

// NewOvertimeRuleRepository creates a new OvertimeRuleRepository.
func NewOvertimeRuleRepository(db *gorm.DB) *OvertimeRuleRepository {
	return &OvertimeRuleRepository{db: db}
}

// --- OvertimeRule CRUD ---

func (r *OvertimeRuleRepository) Create(rule *core.OvertimeRule) error {
	return r.db.Create(rule).Error
}

func (r *OvertimeRuleRepository) Update(rule *core.OvertimeRule) error {
	return r.db.Save(rule).Error
}

func (r *OvertimeRuleRepository) GetByID(id uint) (*core.OvertimeRule, error) {
	var rule core.OvertimeRule
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetActive returns the first active overtime rule.
func (r *OvertimeRuleRepository) GetActive() (*core.OvertimeRule, error) {
	var rule core.OvertimeRule
	err := r.db.Where("is_active = ?", true).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *OvertimeRuleRepository) List() ([]core.OvertimeRule, error) {
	var rules []core.OvertimeRule
	err := r.db.Order("name ASC").Find(&rules).Error
	return rules, err
}

func (r *OvertimeRuleRepository) Delete(id uint) error {
	return r.db.Delete(&core.OvertimeRule{}, id).Error
}

// --- PublicHoliday ---

func (r *OvertimeRuleRepository) CreateHoliday(holiday *core.PublicHoliday) error {
	return r.db.Create(holiday).Error
}

func (r *OvertimeRuleRepository) UpdateHoliday(holiday *core.PublicHoliday) error {
	return r.db.Save(holiday).Error
}

func (r *OvertimeRuleRepository) GetHolidayByID(id uint) (*core.PublicHoliday, error) {
	var holiday core.PublicHoliday
	if err := r.db.First(&holiday, id).Error; err != nil {
		return nil, err
	}
	return &holiday, nil
}

func (r *OvertimeRuleRepository) ListHolidays(year int) ([]core.PublicHoliday, error) {
	var holidays []core.PublicHoliday
	err := r.db.Where("EXTRACT(YEAR FROM date) = ?", year).
		Order("date ASC").Find(&holidays).Error
	return holidays, err
}

func (r *OvertimeRuleRepository) DeleteHoliday(id uint) error {
	return r.db.Delete(&core.PublicHoliday{}, id).Error
}

// IsHoliday checks whether the given date falls on a public holiday.
func (r *OvertimeRuleRepository) IsHoliday(date time.Time) (bool, error) {
	var count int64
	d := date.Truncate(24 * time.Hour)
	err := r.db.Model(&core.PublicHoliday{}).
		Where("date = ?", d).Count(&count).Error
	return count > 0, err
}

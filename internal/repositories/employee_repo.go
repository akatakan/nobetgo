package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *EmployeeRepository) Create(employee *core.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) GetByID(id uint) (*core.Employee, error) {
	var employee core.Employee
	err := r.db.Preload("Department").Preload("Title").First(&employee, id).Error
	return &employee, err
}

func (r *EmployeeRepository) List() ([]core.Employee, error) {
	var employees []core.Employee
	err := r.db.Preload("Department").Preload("Title").Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Update(employee *core.Employee) error {
	return r.db.Save(employee).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&core.Employee{}, id).Error
}

func (r *EmployeeRepository) ListByDepartment(departmentID uint) ([]core.Employee, error) {
	var employees []core.Employee
	err := r.db.Preload("Department").Preload("Title").Where("department_id = ? AND is_active = ?", departmentID, true).Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) GetByEmail(email string) (*core.Employee, error) {
	var employee core.Employee
	err := r.db.Where("email = ?", email).First(&employee).Error
	return &employee, err
}

func (r *EmployeeRepository) ListPaginated(params core.PaginationParams) ([]core.Employee, int64, error) {
	var employees []core.Employee
	var total int64

	db := r.db.Model(&core.Employee{}).Preload("Department").Preload("Title")

	if params.Search != "" {
		search := "%" + params.Search + "%"
		db = db.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", search, search, search)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := db.Offset(offset).Limit(params.Limit).Find(&employees).Error

	return employees, total, err
}

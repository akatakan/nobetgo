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

func (r *EmployeeRepository) Create(employee *core.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) GetByID(id uint) (*core.Employee, error) {
	var employee core.Employee
	err := r.db.First(&employee, id).Error
	return &employee, err
}

func (r *EmployeeRepository) List() ([]core.Employee, error) {
	var employees []core.Employee
	err := r.db.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Update(employee *core.Employee) error {
	return r.db.Save(employee).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&core.Employee{}, id).Error
}

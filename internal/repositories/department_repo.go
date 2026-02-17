package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(department *core.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) GetByID(id uint) (*core.Department, error) {
	var department core.Department
	err := r.db.First(&department, id).Error
	return &department, err
}

func (r *DepartmentRepository) List() ([]core.Department, error) {
	var departments []core.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) Update(department *core.Department) error {
	return r.db.Save(department).Error
}

func (r *DepartmentRepository) Delete(id uint) error {
	return r.db.Delete(&core.Department{}, id).Error
}

func (r *DepartmentRepository) GetByName(name string) (*core.Department, error) {
	var department core.Department
	err := r.db.Where("name = ?", name).First(&department).Error
	return &department, err
}

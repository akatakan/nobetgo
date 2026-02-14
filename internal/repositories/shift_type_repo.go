package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type ShiftTypeRepository struct {
	db *gorm.DB
}

func NewShiftTypeRepository(db *gorm.DB) *ShiftTypeRepository {
	return &ShiftTypeRepository{db: db}
}

func (r *ShiftTypeRepository) Create(shiftType *core.ShiftType) error {
	return r.db.Create(shiftType).Error
}

func (r *ShiftTypeRepository) GetByID(id uint) (*core.ShiftType, error) {
	var shiftType core.ShiftType
	err := r.db.First(&shiftType, id).Error
	return &shiftType, err
}

func (r *ShiftTypeRepository) List() ([]core.ShiftType, error) {
	var shiftTypes []core.ShiftType
	err := r.db.Find(&shiftTypes).Error
	return shiftTypes, err
}

func (r *ShiftTypeRepository) Update(shiftType *core.ShiftType) error {
	return r.db.Save(shiftType).Error
}

func (r *ShiftTypeRepository) Delete(id uint) error {
	return r.db.Delete(&core.ShiftType{}, id).Error
}

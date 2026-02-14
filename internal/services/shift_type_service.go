package services

import (
	"github.com/akatakan/nobetgo/internal/core"
)

type ShiftTypeRepositoryInterface interface {
	Create(shiftType *core.ShiftType) error
	GetByID(id uint) (*core.ShiftType, error)
	List() ([]core.ShiftType, error)
	Update(shiftType *core.ShiftType) error
	Delete(id uint) error
}

type ShiftTypeService struct {
	repo ShiftTypeRepositoryInterface
}

func NewShiftTypeService(repo ShiftTypeRepositoryInterface) *ShiftTypeService {
	return &ShiftTypeService{repo: repo}
}

func (s *ShiftTypeService) CreateShiftType(shiftType *core.ShiftType) error {
	return s.repo.Create(shiftType)
}

func (s *ShiftTypeService) GetShiftTypeByID(id uint) (*core.ShiftType, error) {
	return s.repo.GetByID(id)
}

func (s *ShiftTypeService) GetAllShiftTypes() ([]core.ShiftType, error) {
	return s.repo.List()
}

func (s *ShiftTypeService) UpdateShiftType(shiftType *core.ShiftType) error {
	return s.repo.Update(shiftType)
}

func (s *ShiftTypeService) DeleteShiftType(id uint) error {
	return s.repo.Delete(id)
}

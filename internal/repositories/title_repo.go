package repositories

import (
	"github.com/akatakan/nobetgo/internal/core"
	"gorm.io/gorm"
)

type TitleRepository struct {
	db *gorm.DB
}

func NewTitleRepository(db *gorm.DB) *TitleRepository {
	return &TitleRepository{db: db}
}

func (r *TitleRepository) Create(title *core.Title) error {
	return r.db.Create(title).Error
}

func (r *TitleRepository) GetByID(id uint) (*core.Title, error) {
	var title core.Title
	err := r.db.First(&title, id).Error
	return &title, err
}

func (r *TitleRepository) List() ([]core.Title, error) {
	var titles []core.Title
	err := r.db.Find(&titles).Error
	return titles, err
}

func (r *TitleRepository) Update(title *core.Title) error {
	return r.db.Save(title).Error
}

func (r *TitleRepository) Delete(id uint) error {
	return r.db.Delete(&core.Title{}, id).Error
}

func (r *TitleRepository) GetByName(name string) (*core.Title, error) {
	var title core.Title
	err := r.db.Where("name = ?", name).First(&title).Error
	return &title, err
}

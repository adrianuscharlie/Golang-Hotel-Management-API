package repository

import (
	"hms-backend/model"

	"gorm.io/gorm"
)

type GuestRepository interface {
	Create(m *model.Guest) (*model.Guest, error)
	FindByID(id uint) (*model.Guest, error)
	FindByCredentialID(credType, credID string) (*model.Guest, error)
	Update(m *model.Guest) error
	Delete(credType, credID string) error
}

type guestRepository struct {
	db *gorm.DB
}

func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db}
}

func (r *guestRepository) FindByID(id uint) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.Where("id = ? ", id).First(&guest)
	return &guest, err.Error
}

func (r *guestRepository) Create(m *model.Guest) (*model.Guest, error) {
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *guestRepository) FindByCredentialID(credType, credID string) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.Where("id_number = ? AND credential_type = ?", credID, credType).First(&guest)
	return &guest, err.Error
}

func (r *guestRepository) Update(m *model.Guest) error {
	return r.db.Save(&m).Error
}

func (r *guestRepository) Delete(credType, credID string) error {
	return r.db.Where("id_number = ? and credential_type = ?", credID, credType).Delete(&model.Guest{}).Error
}

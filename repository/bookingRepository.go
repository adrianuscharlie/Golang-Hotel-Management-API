package repository

import (
	"hms-backend/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(b *model.Booking) error
	GetByID(i string) (model.Booking, error)
	Update(b *model.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) Create(b *model.Booking) error {
	return r.db.Table("bookings").Create(b).Error
}

func (r *bookingRepository) GetByID(i string) (model.Booking, error) {
	var booking model.Booking
	err := r.db.
		Preload("Room").
		Preload("Guest").
		First(&booking, "id = ?", i).Error
	return booking, err
}

func (r *bookingRepository) Update(b *model.Booking) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", b.ID).Updates(b).Error

}

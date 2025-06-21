package repository

import (
	"hms-backend/model"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(b *model.Booking) error
	Update(b *model.Booking) error
	FindByID(s string) (*model.Booking, error)
	FindByReferenceID(s string) (*model.Booking, error)
	FindForDateRange(start, end time.Time) ([]*model.Booking, error)
	FindByGuestID(guestID uint) ([]*model.Booking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) Create(b *model.Booking) error {
	return r.db.Create(b).Error
}
func (r *bookingRepository) Update(b *model.Booking) error {
	return r.db.Save(b).Error
}
func (r *bookingRepository) FindByID(s string) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Room").Preload("RoomType").Preload("Guest").Where("id = ?", s).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, err
}
func (r *bookingRepository) FindByReferenceID(s string) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Room").Preload("RoomType").Preload("Guest").Where("booking_reference = ?", s).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, err
}
func (r *bookingRepository) FindForDateRange(start, end time.Time) ([]*model.Booking, error) {
	var bookings []*model.Booking
	err := r.db.Preload("Room").Preload("RoomType").Preload("Guest").Where("check_in_date < ? AND check_out_date > ?", end, start).
		Find(&bookings).Error
	return bookings, err
}
func (r *bookingRepository) FindByGuestID(guestID uint) ([]*model.Booking, error) {
	var bookings []*model.Booking
	err := r.db.Preload("Room").Preload("RoomType").Preload("Guest").Where("guest_id = ?", guestID).Find(&bookings).Error
	return bookings, err
}

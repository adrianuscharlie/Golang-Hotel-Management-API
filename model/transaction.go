package model

import "time"

type Transaction struct {
	Id            uint `gorm:"primaryKey"`
	BookingID     uint
	Booking       Booking
	Amount        float64
	PaymentMethod string
	Paid          bool
	CreatedAt     time.Time
}

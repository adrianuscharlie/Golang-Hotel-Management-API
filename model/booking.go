package model

import "time"

type Booking struct {
	ID           uint `gorm:"primaryKey"`
	RoomID       int
	GuestID      int
	CheckInDate  time.Time
	CheckOutDate time.Time
	Status       string
	CreatedAt    time.Time
	Notes        string
}

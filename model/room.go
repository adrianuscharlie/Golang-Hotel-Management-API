package model

import "time"

// --- Enum-like type for Room Status ---

type RoomStatus string

const (
	// StatusAvailable means the room is in service and can be booked.
	StatusAvailable RoomStatus = "available"
	StatusOccupied  RoomStatus = "occupied"

	// StatusMaintenance means the room is out of service and cannot be booked.
	StatusMaintenance RoomStatus = "maintenance"
)

// --- Room Model ---

type Room struct {
	ID         uint `gorm:"primaryKey"`
	RoomTypeID uint
	Number     string `gorm:"unique;not null"`

	// Use the custom RoomStatus type for type safety.
	// GORM will store its string value (e.g., "available") in the database.
	Status RoomStatus `gorm:"not null;default:available"`

	// Define the relationship for GORM Preload.
	RoomType RoomType

	// Automatic timestamps for auditing.
	CreatedAt time.Time
	UpdatedAt time.Time
}

// --- RoomType Model ---

type RoomType struct {
	ID          uint    `gorm:"primaryKey"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`

	// Name should be unique to prevent duplicate room types.
	Name string `gorm:"unique;not null"`

	// Capacity should be an unsigned integer.
	Capacity uint `gorm:"not null"`
}

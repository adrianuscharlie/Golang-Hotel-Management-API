package model

import "time"

// Define a custom type for BookingStatus for type safety.
type BookingStatus string

const (
	StatusPending    BookingStatus = "pending"
	StatusConfirmed  BookingStatus = "confirmed"
	StatusCheckedIn  BookingStatus = "checked_in"
	StatusCheckedOut BookingStatus = "checked_out"
	StatusCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	// High-performance, unique, time-sortable primary key for internal use.
	ID string `gorm:"primaryKey;type:char(26)"`

	// Human-readable, unique reference for customer-facing communication.
	BookingReference string `gorm:"unique;not null;type:varchar(30)"`

	// --- Foreign Keys and Relationships ---

	// Use uint for foreign keys pointing to auto-incrementing IDs.
	RoomID  uint
	GuestID uint

	// Add the struct fields for GORM relationships. This enables Preload.
	Room  *Room
	Guest *Guest

	// --- Booking Details ---

	CheckInDate  time.Time
	CheckOutDate time.Time

	// Use the custom BookingStatus type to prevent typos.
	// GORM will store this as a string in the database.
	Status BookingStatus `gorm:"type:varchar(20)"`

	// Using `type:text` allows for longer notes if needed.
	Notes string `gorm:"type:text"`

	// --- Automatic Timestamps ---

	// GORM automatically handles these fields by name.
	CreatedAt time.Time
	UpdatedAt time.Time
}

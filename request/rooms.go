package request

import "time"

type CreateRoomRequest struct {
	Number     string `json:"number" binding:"required"`
	Status     string `json:"status" binding:"required"`
	RoomTypeID int    `json:"room_type_id" binding:"required"`
}

type UpdateRoomRequest struct {
	ID         uint   `json:"id" binding:"required"`
	Number     string `json:"number" binding:"required"`
	Status     string `json:"status" binding:"required"`
	RoomTypeID uint   `json:"room_type_id" binding:"required"`
}

type ChangeRoomStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateRoomTypeRequest struct {
	Price       float64 `json:"price"`
	Capacity    uint    `json:"capacity"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
}

type RoomFilterParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Category string
	MinPrice float64
	MaxPrice float64
}

type ChangeStatus struct {
	ID     string `json:"room_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

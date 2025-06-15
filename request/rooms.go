package request

type CreateRoomRequest struct {
	Number     string `json:"number" binding:"required"`
	Status     string `json:"status" binding:"required"`
	RoomTypeID int    `json:"room_type_id" binding:"required"`
}

type UpdateRoomRequest struct {
	ID         int    `json:"id" binding:"required"`
	Number     string `json:"number" binding:"required"`
	Status     string `json:"status" binding:"required"`
	RoomTypeID int    `json:"room_type_id" binding:"required"`
}

type ChangeRoomStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateRoomTypeRequest struct {
	Price       float64 `json:"price"`
	Capacity    int     `json:"capacity"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
}

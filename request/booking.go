package request

type CreateBookingRequest struct {
	RoomID       uint   `json:"room_id" binding:"required"`
	GuestID      uint   `json:"guest_id" binding:"required"`
	CheckInDate  string `json:"check_in_date" binding:"required"`
	CheckOutDate string `json:"check_out_date" binding:"required"`
	Notes        string `json:"notes"`
}

type UpdateBookingRequest struct {
}

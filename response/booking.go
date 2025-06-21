package response

import "hms-backend/model"

type BookingResponse struct {
	BookingID      string                              `json:"booking_id" binding:"required"`
	CheckInDate    string                              `json:"check_in_date" binding:"required"`
	CheckOutDate   string                              `json:"check_out_date" binding:"required"`
	Status         model.BookingStatus                 `json:"status"`
	Notes          string                              `json:"notes"`
	AdditionalInfo AdditionalInfoCreateBookingResponse `json:"additionalInfo"`
}

type AdditionalInfoCreateBookingResponse struct {
	Room  RoomResponse  `json:"room"`
	Guest GuestResponse `json:"guest"`
}

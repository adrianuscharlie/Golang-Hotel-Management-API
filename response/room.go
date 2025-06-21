package response

import "hms-backend/model"

type RoomResponse struct {
	ID       uint             `json:"id"`
	Number   string           `json:"number"`
	Status   model.RoomStatus `json:"status"`
	RoomType RoomTypeDetail   `json:"room_type"`
}

type RoomTypeDetail struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    uint    `json:"capacity"`
	Price       float64 `json:"price"`
}

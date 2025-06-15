package response

type RoomResponse struct {
	ID       int            `json:"id"`
	Number   string         `json:"number"`
	Status   string         `json:"status"`
	RoomType RoomTypeDetail `json:"room_type"`
}

type RoomTypeDetail struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
}

package handler

import (
	"hms-backend/request"
	"hms-backend/response"
	"hms-backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomServices services.RoomServices
}

// Constructor
func NewRoomHandler(s services.RoomServices) *RoomHandler {
	return &RoomHandler{roomServices: s}
}

// GET /rooms
func (h *RoomHandler) GetAllRoom(c *gin.Context) {
	resp, err := h.roomServices.GetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, response.Response{"400", "Room Not Found: " + err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Successful", resp})
}

// POST /rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req request.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", "Invalid request: " + err.Error(), nil})
		return
	}

	// pass by pointer since service expects *CreateRoomRequest
	room, err := h.roomServices.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", "Failed to create room: " + err.Error(), nil})
		return
	}
	c.JSON(http.StatusCreated, response.Response{"00", "Successful creating new room", room})
}

// POST /rooms/type
func (h *RoomHandler) CreateRoomType(c *gin.Context) {
	var req request.CreateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}
	err := h.roomServices.CreateRoomType(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusCreated, response.Response{"00", "Successful", nil})
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id := c.Param("id")
	roomID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}
	// Use roomID to query from DB, for example
	room, err := h.roomServices.GetByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Response{"404", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", room})
}

func (h *RoomHandler) GetAvailableRoom(c *gin.Context) {

	checkInQuery := c.Query("check_in")
	checkOutQuery := c.Query("check_out")
	layout := "2006-01-02"
	if checkInQuery == "" || checkOutQuery == "" {
		c.JSON(http.StatusBadRequest, response.Response{"400", "check_in and check_out dates are required", nil})
		return
	}

	checkInStr, err := time.Parse(layout, checkInQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", "invalid date check_in format", nil})
		return
	}

	checkOutStr, err := time.Parse(layout, checkOutQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", "invalid date check_out format", nil})
		return
	}

	params := request.RoomFilterParams{
		CheckIn:  checkInStr,
		CheckOut: checkOutStr,
		Category: c.Query("category"),
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		params.MinPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		params.MaxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	rooms, err := h.roomServices.FindAvailable(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"200", "Successful", rooms})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	var req request.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}
	room, err := h.roomServices.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful update room", room})
}

func (h *RoomHandler) ChangeStatus(c *gin.Context) {
	var req request.ChangeStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}

	strId, err := strconv.Atoi(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}
	err = h.roomServices.ChangeStatus(uint(strId), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Successful", nil})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	room := c.Query("room_number")
	if room == "" {
		c.JSON(http.StatusBadRequest, response.Response{"400", "Room number required", nil})
		return
	}
	if err := h.roomServices.Delete(room); err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Successful delete room", nil})

}

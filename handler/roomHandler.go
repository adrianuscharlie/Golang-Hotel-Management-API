package handler

import (
	"hms-backend/request"
	"hms-backend/services"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusNotFound, gin.H{
			"responseCode":    "404",
			"responseMessage": "Error Get All Room: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "00",
		"responseMessage": "Success",
		"data": gin.H{
			"rooms": resp,
		},
	})
}

// POST /rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req request.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"responseCode":    "400",
			"responseMessage": "Invalid request: " + err.Error(),
		})
		return
	}

	// pass by pointer since service expects *CreateRoomRequest
	if err := h.roomServices.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"responseCode":    "500",
			"responseMessage": "Failed to create room: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"responseCode":    "00",
		"responseMessage": "Room successfully created",
	})
}

// POST /rooms/type
func (h *RoomHandler) CreateRoomType(c *gin.Context) {
	var req request.CreateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"responseCode":    "400",
			"responseMessage": err.Error(),
		})
		return
	}
	err := h.roomServices.CreateRoomType(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"responseCode":    "500",
			"responseMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"responseCode":    "00",
		"responseMessage": "Success",
	})
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id := c.Param("id")
	roomID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"responseCode":    "400",
			"responseMessage": err.Error(),
		})
		return
	}
	// Use roomID to query from DB, for example
	room, err := h.roomServices.GetByID(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"responseCode":    "404",
			"responseMessage": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "00",
		"responseMessage": "Success",
		"data":            room,
	})
}

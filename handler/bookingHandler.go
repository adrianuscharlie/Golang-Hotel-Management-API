package handler

import (
	"hms-backend/request"
	"hms-backend/response"
	"hms-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService services.BookingServices
}

func NewBookingHandler(s services.BookingServices) *BookingHandler {
	return &BookingHandler{bookingService: s}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}

	res, err := h.bookingService.CreateBooking(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusCreated, response.Response{"00", "Sucessful", res})
}

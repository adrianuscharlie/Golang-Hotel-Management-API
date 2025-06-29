package handler

import (
	"hms-backend/request"
	"hms-backend/response"
	"hms-backend/services"
	"net/http"
	"time"

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

func (h *BookingHandler) GetBookingByReference(c *gin.Context) {
	ref := c.Param("id")
	if ref == "" {
		c.JSON(http.StatusNotFound, response.Response{"404", "Bad Request", nil})
		return
	}
	res, err := h.bookingService.GetBookingByReference(ref)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}

func (h *BookingHandler) GetBookingByDateRange(c *gin.Context) {
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

	bookings, err := h.bookingService.ListBookingsForDateRange(checkInStr, checkOutStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", bookings})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	var req *request.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}
	res, err := h.bookingService.CancelBooking(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}

func (h *BookingHandler) CheckIn(c *gin.Context) {
	var req request.CheckInCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, response.Response{"404", "Bad Request", nil})
		return
	}
	res, err := h.bookingService.CheckInGuest(req.BookingReference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}
func (h *BookingHandler) Checkout(c *gin.Context) {
	var req request.CheckInCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, response.Response{"404", "Bad Request", nil})
		return
	}
	res, err := h.bookingService.CheckOutGuest(req.BookingReference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}

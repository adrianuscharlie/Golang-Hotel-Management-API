package handler

import (
	"hms-backend/request"
	"hms-backend/response"
	"hms-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GuestHandler struct {
	guestServices services.GuestService
}

func NewGuestHandler(s services.GuestService) *GuestHandler {
	return &GuestHandler{guestServices: s}
}

func (h *GuestHandler) CreateNewGuest(c *gin.Context) {
	var req request.GuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}

	res, err := h.guestServices.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusCreated, response.Response{"00", "Sucessful", res})
}

func (h *GuestHandler) GetGuestByCredentialID(c *gin.Context) {
	credType := c.Query("credential_type")
	credID := c.Query("credential_number")
	res, err := h.guestServices.FindByCredentialID(credType, credID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Response{"404", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}

func (h *GuestHandler) Update(c *gin.Context) {
	var req request.GuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{"400", err.Error(), nil})
		return
	}

	res, err := h.guestServices.Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{"500", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", res})
}

func (h *GuestHandler) Delete(c *gin.Context) {
	credType := c.Query("credential_type")
	credID := c.Query("credential_number")
	err := h.guestServices.Delete(credType, credID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Response{"404", err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, response.Response{"00", "Sucessful", nil})
}

package routes

import (
	"hms-backend/handler"
	"hms-backend/repository"
	"hms-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize Repositories, Services, Handlers
	roomRepository := repository.NewRoomRepository(db)
	roomServices := services.NewRoomServices(roomRepository)
	roomHandler := handler.NewRoomHandler(roomServices)

	guestRepository := repository.NewGuestRepository(db)
	guestServices := services.NewGuestServices(guestRepository)
	guestHandler := handler.NewGuestHandler(guestServices)

	bookingRepository := repository.NewBookingRepository(db)
	bookingServices := services.NewBookingServices(bookingRepository, roomServices, guestServices)
	bookingHandler := handler.NewBookingHandler(bookingServices)
	// Main API group
	api := router.Group("/api")
	{
		// Room routes group: /api/room
		roomApi := api.Group("/room")
		{
			roomApi.POST("/", roomHandler.CreateRoom)         // POST /api/room
			roomApi.POST("/type", roomHandler.CreateRoomType) // POST /api/room/type
			roomApi.GET("/", roomHandler.GetAllRoom)          // GET  /api/room
			roomApi.GET("/:id", roomHandler.GetRoomByID)      // GET  /api/room/:id  <-- added
			roomApi.GET("/available", roomHandler.GetAvailableRoom)
			roomApi.PUT("/", roomHandler.UpdateRoom)
			roomApi.PUT("/status", roomHandler.ChangeStatus)
			roomApi.DELETE("/", roomHandler.DeleteRoom)
		}

		guestApi := api.Group("/guest")
		{
			guestApi.POST("/", guestHandler.CreateNewGuest)
			guestApi.GET("/identity", guestHandler.GetGuestByCredentialID)
			guestApi.PUT("/", guestHandler.Update)
			guestApi.DELETE("/:id", guestHandler.Delete)
		}

		bookingApi := api.Group("/booking")
		{
			bookingApi.POST("/", bookingHandler.CreateBooking)
		}

		// You can add other groups here, like:
		// guestApi := api.Group("/guest")
		// bookingApi := api.Group("/booking")
	}
}

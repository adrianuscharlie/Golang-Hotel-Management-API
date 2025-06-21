package services

import (
	"crypto/rand"
	"fmt"
	"hms-backend/model"
	"hms-backend/repository"
	"hms-backend/request"
	"hms-backend/response"
	"time"

	"github.com/oklog/ulid/v2"
)

type BookingServices interface {
	CreateBooking(req *request.CreateBookingRequest) (*response.BookingResponse, error)
	// GetBookingByReference(ref string) (*response.BookingResponse, error)
	// ListBookingsForGuest(id uint) ([]*response.BookingResponse, error)
	// ListBookingsForDateRange(start, end time.Time) ([]*response.BookingResponse, error)
	// CancelBooking(ref string) (*response.BookingResponse, error)
	// CheckInGuest(ref string) (*response.BookingResponse, error)
	// CheckOutGuest(ref string) (*response.BookingResponse, error)
}

type bookingService struct {
	bookingRepository repository.BookingRepository
	roomServices      RoomServices
	guestServices     GuestService
}

func NewBookingServices(repo repository.BookingRepository, room RoomServices, guest GuestService) BookingServices {
	return &bookingService{bookingRepository: repo, roomServices: room, guestServices: guest}
}

func (s *bookingService) CreateBooking(req *request.CreateBookingRequest) (*response.BookingResponse, error) {
	layout := "2006-01-02"
	room, err := s.roomServices.GetRoomModelByID(req.RoomID)
	if err != nil {
		return nil, err
	}
	guest, err := s.guestServices.FindByModelID(req.GuestID)
	if err != nil {
		return nil, err
	}
	ref, err := generateBookingReference()
	if err != nil {
		return nil, err
	}
	checkInStr, err := time.Parse(layout, req.CheckInDate)
	if err != nil {
		return nil, err
	}
	checkoutStr, err := time.Parse(layout, req.CheckOutDate)
	if err != nil {
		return nil, err
	}

	newBooking := model.Booking{
		ID:               generateULID(),
		BookingReference: ref,
		RoomID:           req.RoomID,
		GuestID:          req.GuestID,
		Room:             room,
		Guest:            guest,
		CheckInDate:      checkInStr,
		CheckOutDate:     checkoutStr,
		Status:           model.StatusConfirmed,
		Notes:            req.Notes,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err = s.bookingRepository.Create(&newBooking)
	if err != nil {
		return nil, err
	}
	roomResponse, err := s.roomServices.GetByID(req.RoomID)
	if err != nil {
		return nil, err
	}
	guestResponse, err := s.guestServices.FindByID(req.GuestID)
	if err != nil {
		return nil, err
	}
	return &response.BookingResponse{
		BookingID:    ref,
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate,
		Status:       newBooking.Status,
		Notes:        newBooking.Notes,
		AdditionalInfo: response.AdditionalInfoCreateBookingResponse{
			Room:  *roomResponse,
			Guest: *guestResponse,
		},
	}, err
}

const (
	idCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	idLength  = 4
)

func generateULID() string {
	// Note: In a high-concurrency production app, you might want to initialize
	// the entropy source once and share it, but for most uses, this is fine.
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
func generateBookingReference() (string, error) {
	// 1. Get the prefix and date
	prefix := "BK"
	date := time.Now().Format("20060102") // YYYYMMDD format

	// 2. Generate the random part
	randomPart := make([]byte, idLength)
	if _, err := rand.Read(randomPart); err != nil {
		return "", err
	}
	for i := 0; i < idLength; i++ {
		randomPart[i] = idCharset[int(randomPart[i])%len(idCharset)]
	}

	// 3. Combine them
	return fmt.Sprintf("%s-%s-%s", prefix, date, string(randomPart)), nil
}

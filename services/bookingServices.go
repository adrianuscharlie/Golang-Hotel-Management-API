package services

import (
	"crypto/rand"
	"errors"
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
	GetBookingByReference(ref string) (*response.BookingResponse, error)
	ListBookingsForGuest(id uint) ([]*response.BookingResponse, error)
	ListBookingsForDateRange(start, end time.Time) ([]*response.BookingResponse, error)
	CancelBooking(req *request.CancelBookingRequest) (*response.BookingResponse, error)
	CheckInGuest(ref string) (*response.BookingResponse, error)
	CheckOutGuest(ref string) (*response.BookingResponse, error)
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
	if room.Status != model.StatusAvailable {
		return nil, errors.New("room not available, please use another room")
	}
	guest, err := s.guestServices.FindByModelID(req.GuestID)
	if err != nil {
		return nil, err
	}
	book, err := s.bookingRepository.FindByGuestID(req.GuestID)
	if len(book) > 0 || err != nil {
		return nil, errors.New("guest already have another booking" + err.Error())
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
		Status:           model.StatusPending,
		Notes:            req.Notes,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err = s.bookingRepository.Create(&newBooking)
	if err != nil {
		return nil, err
	}

	return mapToBookingResponse(&newBooking), err
}

func (s *bookingService) GetBookingByReference(ref string) (*response.BookingResponse, error) {
	book, err := s.bookingRepository.FindByReferenceID(ref)
	if err != nil {
		return nil, err
	}
	return mapToBookingResponse(book), err
}

func (s *bookingService) ListBookingsForGuest(id uint) ([]*response.BookingResponse, error) {
	bookings, err := s.bookingRepository.FindByGuestID(id)
	if err != nil {
		return nil, err
	}
	return mapToBookingResponseSlice(bookings), err

}

func (s *bookingService) CancelBooking(req *request.CancelBookingRequest) (*response.BookingResponse, error) {
	booking, err := s.bookingRepository.FindByReferenceID(req.BookingReference)
	if err != nil {
		return nil, errors.New("Booking Not Found")
	}
	booking.Status = model.StatusCancelled
	booking.Notes = req.Reason
	booking.UpdatedAt = time.Now()
	err = s.bookingRepository.Update(booking)
	if err != nil {
		return nil, errors.New("Failed To Cancel Booking")
	}
	return mapToBookingResponse(booking), err
}

func (s *bookingService) ListBookingsForDateRange(start, end time.Time) ([]*response.BookingResponse, error) {
	bookings, err := s.bookingRepository.FindForDateRange(start, end)
	if err != nil {
		return nil, err
	}
	return mapToBookingResponseSlice(bookings), err
}

func (s *bookingService) CheckInGuest(ref string) (*response.BookingResponse, error) {
	booking, err := s.bookingRepository.FindByReferenceID(ref)
	if err != nil {
		return nil, errors.New("Booking Not Found")
	}
	if booking.Status != model.StatusConfirmed {
		return nil, errors.New("Booking Status Not Confirmed")
	}
	if booking.CheckInDate.Day() != time.Now().Day() {
		return nil, errors.New("Cannot Checkin Before/After Day Checkin")
	}

	booking.Status = model.StatusCheckedIn
	booking.UpdatedAt = time.Now()
	err = s.bookingRepository.Update(booking)
	if err != nil {
		return nil, errors.New("Checkin Failed!")
	}
	return mapToBookingResponse(booking), nil
}
func (s *bookingService) CheckOutGuest(ref string) (*response.BookingResponse, error) {
	booking, err := s.bookingRepository.FindByReferenceID(ref)
	if err != nil {
		return nil, errors.New("Booking Not Found")
	}
	if booking.Status != model.StatusConfirmed {
		return nil, errors.New("Booking Status Not Confirmed")
	}
	if booking.CheckOutDate.Day() > time.Now().Day() {
		booking.Notes = "Late Checkout, Must Be Charged for Extra"
	}
	booking.Status = model.StatusCheckedOut
	booking.UpdatedAt = time.Now()
	err = s.bookingRepository.Update(booking)
	if err != nil {
		return nil, errors.New("Checkout Failed!")
	}
	err = s.roomServices.ChangeStatus(booking.RoomID, string(model.StatusAvailable))
	if err != nil {
		//should do something
	}
	return mapToBookingResponse(booking), nil
}

func mapToBookingResponse(booking *model.Booking) *response.BookingResponse {
	if booking == nil {
		return nil
	}
	layout := "2006-01-02"
	resp := &response.BookingResponse{
		BookingID:    booking.BookingReference,
		CheckInDate:  booking.CheckInDate.Format(layout),
		CheckOutDate: booking.CheckOutDate.Format(layout),
		Status:       booking.Status,
		Notes:        booking.Notes,
	}
	if booking.Room != nil {
		resp.AdditionalInfo.Room = *mapToRoomDetail(booking.Room)
	}
	if booking.Guest != nil {
		resp.AdditionalInfo.Guest = *mapToGuestDetail(booking.Guest)
	}
	return resp
}

func mapToRoomDetail(room *model.Room) *response.RoomResponse {
	if room == nil {
		return nil
	}
	resp := &response.RoomResponse{
		ID:     room.ID,
		Number: room.Number,
		Status: room.Status,
	}
	if &room.RoomType != nil {
		resp.RoomType = response.RoomTypeDetail{
			ID:          room.RoomType.ID,
			Name:        room.RoomType.Name,
			Description: room.RoomType.Description,
			Capacity:    room.RoomType.Capacity,
			Price:       room.RoomType.Price,
		}
	}
	return resp
}

func mapToGuestDetail(guest *model.Guest) *response.GuestResponse {
	if guest == nil {
		return nil
	}
	return &response.GuestResponse{
		CredentialType: guest.CredentialType,
		IDNumber:       guest.IDNumber,
		FullName:       guest.FullName,
		Email:          guest.Email,
		Phone:          guest.Phone,
		ID:             guest.ID,
	}
}

func mapToBookingResponseSlice(bookings []*model.Booking) []*response.BookingResponse {
	bookingResponse := make([]*response.BookingResponse, len(bookings))
	for i, booking := range bookings {
		bookingResponse[i] = mapToBookingResponse(booking)
	}
	return bookingResponse
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

package services

import (
	"errors"
	"fmt"
	"hms-backend/model"
	"hms-backend/repository"
	"hms-backend/request"
	"hms-backend/response"
)

type RoomServices interface {
	GetAll() ([]response.RoomResponse, error)
	GetByID(id uint) (*response.RoomResponse, error)
	Create(input *request.CreateRoomRequest) (*model.Room, error)
	Update(input request.UpdateRoomRequest) (*model.Room, error)
	Delete(roomNumber string) error
	ChangeStatus(id uint, status string) error
	FindAvailable(params request.RoomFilterParams) ([]response.RoomResponse, error)
	CreateRoomType(input *request.CreateRoomTypeRequest) error
	GetRoomModelByID(id uint) (*model.Room, error)
}

type roomServices struct {
	roomRepository repository.RoomRepository
}

func NewRoomServices(roomRepo repository.RoomRepository) RoomServices {
	return &roomServices{roomRepo}
}

func (s *roomServices) GetAll() ([]response.RoomResponse, error) {
	rooms, err := s.roomRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return mapToRoomResponseSlice(rooms), nil
}

func (s *roomServices) GetByID(id uint) (*response.RoomResponse, error) {
	room, err := s.roomRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return mapToRoomResponse(room), err
}

func (s *roomServices) GetRoomModelByID(id uint) (*model.Room, error) {
	return s.roomRepository.FindByID(id)
}

func (s *roomServices) Create(input *request.CreateRoomRequest) (*model.Room, error) {
	if !isValidRoomStatus(input.Status) {
		return nil, errors.New("invalid room status provided. must be 'available' or 'maintenance'")
	}
	room := model.Room{
		Status:     model.RoomStatus(input.Status),
		Number:     input.Number,
		RoomTypeID: uint(input.RoomTypeID),
	}
	err := s.roomRepository.Create(&room)
	if err != nil {
		return nil, err
	}
	return &room, err
}

func (s *roomServices) Update(update request.UpdateRoomRequest) (*model.Room, error) {
	room, err := s.roomRepository.FindByID(update.ID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}
	if !isValidRoomStatus(update.Status) {
		return nil, errors.New("invalid room status provided. must be 'available' or 'maintenance'")
	}
	room.Number = update.Number
	room.Status = model.RoomStatus(update.Status)
	room.RoomTypeID = uint(update.RoomTypeID)
	room.RoomType = model.RoomType{}
	err = s.roomRepository.Update(room)
	if err != nil {
		return nil, err
	}
	updatedRoom, err := s.roomRepository.FindByID(update.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room after update: %w", err)
	}
	return updatedRoom, nil
}
func (s *roomServices) Delete(roomNumber string) error {
	room, err := s.roomRepository.FindByNumber(roomNumber)
	if err != nil {
		return err
	}
	return s.roomRepository.Delete(room.ID)
}

func (s *roomServices) ChangeStatus(id uint, status string) error {
	if !isValidRoomStatus(status) {
		return errors.New("invalid room status provided. must be 'available' or 'maintenance'")
	}
	return s.roomRepository.ChangeStatus(id, status)
}

func (s *roomServices) FindAvailable(params request.RoomFilterParams) ([]response.RoomResponse, error) {
	rooms, err := s.roomRepository.FindAvailable(params)
	if err != nil {
		return nil, err
	}

	return mapToRoomResponseSlice(rooms), nil
}

func (s *roomServices) CreateRoomType(input *request.CreateRoomTypeRequest) error {
	roomType := model.RoomType{
		Price:       input.Price,
		Capacity:    input.Capacity,
		Description: input.Description,
		Name:        input.Name,
	}
	return s.roomRepository.CreateRoomType(&roomType)
}

func isValidRoomStatus(status string) bool {
	// Cast the string to a RoomStatus to compare against the constants
	s := model.RoomStatus(status)
	switch s {
	case model.StatusAvailable, model.StatusMaintenance:
		return true
	default:
		return false
	}
}

func mapToRoomResponse(room *model.Room) *response.RoomResponse {
	if room == nil {
		return nil
	}
	resp := &response.RoomResponse{
		ID:     room.ID,
		Number: room.Number,
		Status: room.Status, // Correct type cast
	}
	// FIX: Nil check before accessing nested struct fields
	if &room.RoomType != nil {
		resp.RoomType = response.RoomTypeDetail{
			ID:          room.RoomType.ID,
			Name:        room.RoomType.Name,
			Description: room.RoomType.Description,
			Price:       room.RoomType.Price,
			Capacity:    room.RoomType.Capacity,
		}
	}
	return resp
}

func mapToRoomResponseSlice(rooms []model.Room) []response.RoomResponse {
	roomResponses := make([]response.RoomResponse, len(rooms))
	for i, room := range rooms {
		// Since mapToRoomResponse returns a pointer, we dereference it here.
		roomResponses[i] = *mapToRoomResponse(&room)
	}
	return roomResponses
}

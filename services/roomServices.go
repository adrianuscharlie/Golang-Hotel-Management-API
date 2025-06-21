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
	GetByID(id int) (*response.RoomResponse, error)
	Create(input *request.CreateRoomRequest) (*model.Room, error)
	Update(input request.UpdateRoomRequest) (*model.Room, error)
	Delete(roomNumber string) error
	ChangeStatus(id int, status string) error
	FindAvailable(params request.RoomFilterParams) ([]response.RoomResponse, error)
	CreateRoomType(input *request.CreateRoomTypeRequest) error
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
	var roomResponses []response.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, response.RoomResponse{
			ID:     int(room.ID),
			Number: room.Number,
			Status: room.Status,
			RoomType: response.RoomTypeDetail{
				ID:          int(room.RoomType.ID),
				Name:        room.RoomType.Name,
				Description: room.RoomType.Description,
				Price:       room.RoomType.Price,
				Capacity:    room.RoomType.Capacity,
			},
		})
	}
	return roomResponses, nil
}

func (s *roomServices) GetByID(id int) (*response.RoomResponse, error) {
	room, err := s.roomRepository.FindByID(id)
	if err != nil {
		return &response.RoomResponse{}, err
	}
	response := response.RoomResponse{
		ID:     int(room.ID),
		Number: room.Number,
		Status: room.Status,
		RoomType: response.RoomTypeDetail{
			ID:          int(room.RoomType.ID),
			Name:        room.RoomType.Name,
			Description: room.RoomType.Description,
			Price:       room.RoomType.Price,
			Capacity:    room.RoomType.Capacity,
		},
	}
	return &response, err
}

func (s *roomServices) Create(input *request.CreateRoomRequest) (*model.Room, error) {
	room := model.Room{
		Status:     input.Status,
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
	room.Number = update.Number
	room.Status = update.Status
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
	return s.roomRepository.Delete(int(room.ID))
}

func (s *roomServices) ChangeStatus(id int, status string) error {
	return s.roomRepository.ChangeStatus(id, status)
}

func (s *roomServices) FindAvailable(params request.RoomFilterParams) ([]response.RoomResponse, error) {
	rooms, err := s.roomRepository.FindAvailable(params)
	if err != nil {
		return nil, err
	}
	var roomResponses []response.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, response.RoomResponse{
			ID:     int(room.ID),
			Number: room.Number,
			Status: room.Status,
			RoomType: response.RoomTypeDetail{
				ID:          int(room.RoomType.ID),
				Name:        room.RoomType.Name,
				Description: room.RoomType.Description,
				Price:       room.RoomType.Price,
				Capacity:    room.RoomType.Capacity,
			},
		})
	}
	return roomResponses, nil
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

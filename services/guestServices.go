package services

import (
	"hms-backend/model"
	"hms-backend/repository"
	"hms-backend/request"
	"hms-backend/response"
	"time"
)

type GuestService interface {
	Create(r *request.GuestRequest) (*response.GuestResponse, error)
	FindByCredentialID(credType, credID string) (*response.GuestResponse, error)
	FindByID(id uint) (*response.GuestResponse, error)
	FindByModelID(id uint) (*model.Guest, error)
	Update(r *request.GuestRequest) (*response.GuestResponse, error)
	Delete(credType, credID string) error
}

type guestService struct {
	guestRepository repository.GuestRepository
}

func NewGuestServices(repo repository.GuestRepository) GuestService {
	return &guestService{guestRepository: repo}
}

func (s *guestService) FindByModelID(id uint) (*model.Guest, error) {
	return s.guestRepository.FindByID(id)
}

func (s *guestService) FindByID(id uint) (*response.GuestResponse, error) {
	guest, err := s.guestRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return mapGuestResponse(guest), err
}

func (s *guestService) Create(r *request.GuestRequest) (*response.GuestResponse, error) {
	guest := model.Guest{
		CredentialType: r.CredentialType,
		IDNumber:       r.IDNumber,
		FullName:       r.FullName,
		Email:          r.Email,
		Phone:          r.Phone,
		CreatedAt:      time.Now(),
	}
	newGuest, err := s.guestRepository.Create(&guest)
	if err != nil {
		return nil, err
	}
	return mapGuestResponse(newGuest), err
}

func (s *guestService) FindByCredentialID(credType, credID string) (*response.GuestResponse, error) {
	guest, err := s.guestRepository.FindByCredentialID(credType, credID)
	if err != nil {
		return nil, err
	}
	return &response.GuestResponse{
		IDNumber:       guest.IDNumber,
		CredentialType: guest.CredentialType,
		FullName:       guest.FullName,
		Email:          guest.Email,
		Phone:          guest.Phone,
	}, err
}

func (s *guestService) Update(r *request.GuestRequest) (*response.GuestResponse, error) {
	guest, err := s.guestRepository.FindByCredentialID(r.CredentialType, r.IDNumber)
	if err != nil {
		return nil, err
	}
	guest.FullName = r.FullName
	guest.Email = r.Email
	guest.Phone = r.Phone
	guest.IDNumber = r.IDNumber
	err = s.guestRepository.Update(guest)
	if err != nil {
		return nil, err
	}
	return &response.GuestResponse{
		IDNumber:       guest.IDNumber,
		CredentialType: guest.CredentialType,
		FullName:       guest.FullName,
		Email:          guest.Email,
		Phone:          guest.Phone,
	}, err
}

func (s *guestService) Delete(credType, credID string) error {
	return s.guestRepository.Delete(credType, credID)
}

func mapGuestResponse(g *model.Guest) *response.GuestResponse {
	return &response.GuestResponse{
		IDNumber:       g.IDNumber,
		CredentialType: g.CredentialType,
		FullName:       g.FullName,
		Email:          g.Email,
		Phone:          g.Phone,
		ID:             g.ID,
	}
}

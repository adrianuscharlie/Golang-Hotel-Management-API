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
	Update(r *request.GuestRequest) (*response.GuestResponse, error)
	Delete(credType, credID string) error
}

type guestService struct {
	guestRepository repository.GuestRepository
}

func NewGuestServices(repo repository.GuestRepository) GuestService {
	return &guestService{guestRepository: repo}
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
	err := s.guestRepository.Create(&guest)
	if err != nil {
		return nil, err
	}
	res := response.GuestResponse{
		IDNumber:       guest.IDNumber,
		CredentialType: guest.CredentialType,
		FullName:       guest.FullName,
		Email:          guest.Email,
		Phone:          guest.Phone,
	}
	return &res, err
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

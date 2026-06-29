package service

import (
	"errors"

	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	CancelReservation(id uint, userID uint) error
	GetAllReservations() ([]dto.ReservationResponse, error)
}

type reservationServiceImpl struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationServiceImpl{reservationRepo: reservationRepo}
}

func (s *reservationServiceImpl) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	err := s.reservationRepo.CreateReservationSafe(reservation)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationServiceImpl) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.GetReservationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ReservationResponse
	for _, r := range reservations {
		resp := dto.ReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		}
		if r.Zone != nil {
			resp.Zone = &dto.ZoneResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			}
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *reservationServiceImpl) CancelReservation(id uint, userID uint) error {
	reservation, err := s.reservationRepo.GetReservationByIDAndUser(id, userID)
	if err != nil {
		return errors.New("reservation not found or you don't have permission")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	return s.reservationRepo.UpdateReservationStatus(id, "cancelled")
}

func (s *reservationServiceImpl) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var responses []dto.ReservationResponse
	for _, r := range reservations {
		resp := dto.ReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		}
		if r.User != nil {
			resp.User = &dto.UserResponse{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			}
		}
		if r.Zone != nil {
			resp.Zone = &dto.ZoneResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			}
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

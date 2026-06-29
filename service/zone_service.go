package service

import (
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneServiceImpl struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneServiceImpl{zoneRepo: zoneRepo}
}

func (s *zoneServiceImpl) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	err := s.zoneRepo.CreateZone(zone)
	if err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity, // Newly created zone has full capacity available
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *zoneServiceImpl) GetAllZones() ([]dto.ZoneResponse, error) {
	zonesWithAvail, err := s.zoneRepo.GetAllZones()
	if err != nil {
		return nil, err
	}

	var responses []dto.ZoneResponse
	for _, z := range zonesWithAvail {
		responses = append(responses, dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: z.AvailableSpots,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
			UpdatedAt:      z.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *zoneServiceImpl) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	z, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: z.AvailableSpots,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
		UpdatedAt:      z.UpdatedAt,
	}, nil
}

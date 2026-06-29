package repository

import (
	"spotsync/models"

	"gorm.io/gorm"
)

type ZoneWithAvailability struct {
	models.ParkingZone
	AvailableSpots int `json:"available_spots"`
}

type ZoneRepository interface {
	CreateZone(zone *models.ParkingZone) error
	GetAllZones() ([]ZoneWithAvailability, error)
	GetZoneByID(id uint) (*ZoneWithAvailability, error)
}

type zoneRepositoryImpl struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepositoryImpl{db: db}
}

func (r *zoneRepositoryImpl) CreateZone(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepositoryImpl) GetAllZones() ([]ZoneWithAvailability, error) {
	var zones []ZoneWithAvailability
	
	err := r.db.Model(&models.ParkingZone{}).
		Select("parking_zones.*, (parking_zones.total_capacity - (SELECT COUNT(*) FROM reservations WHERE reservations.zone_id = parking_zones.id AND reservations.status = 'active')) AS available_spots").
		Scan(&zones).Error

	return zones, err
}

func (r *zoneRepositoryImpl) GetZoneByID(id uint) (*ZoneWithAvailability, error) {
	var zone ZoneWithAvailability
	
	err := r.db.Model(&models.ParkingZone{}).
		Select("parking_zones.*, (parking_zones.total_capacity - (SELECT COUNT(*) FROM reservations WHERE reservations.zone_id = parking_zones.id AND reservations.status = 'active')) AS available_spots").
		Where("parking_zones.id = ?", id).
		Scan(&zone).Error

	if err != nil {
		return nil, err
	}
	if zone.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &zone, nil
}

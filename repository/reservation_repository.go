package repository

import (
	"errors"

	"spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateReservationSafe(reservation *models.Reservation) error
	GetReservationsByUserID(userID uint) ([]models.Reservation, error)
	GetReservationByIDAndUser(id uint, userID uint) (*models.Reservation, error)
	UpdateReservationStatus(id uint, status string) error
	GetAllReservations() ([]models.Reservation, error)
}

type reservationRepositoryImpl struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepositoryImpl{db: db}
}

// CreateReservationSafe uses a database transaction and row-level locking
// to safely prevent the "EV Spot Bottleneck" race condition.
func (r *reservationRepositoryImpl) CreateReservationSafe(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneID).Error; err != nil {
			return err
		}


		var activeCount int64
		if err := tx.Model(&models.Reservation{}).Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").Count(&activeCount).Error; err != nil {
			return err
		}


		if int(activeCount) >= zone.TotalCapacity {
			return errors.New("zone is at full capacity")
		}


		if err := tx.Create(reservation).Error; err != nil {
			return err
		}


		reservation.Zone = &zone

		return nil
	})
}

func (r *reservationRepositoryImpl) GetReservationsByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepositoryImpl) GetReservationByIDAndUser(id uint, userID uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepositoryImpl) UpdateReservationStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *reservationRepositoryImpl) GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	return reservations, err
}

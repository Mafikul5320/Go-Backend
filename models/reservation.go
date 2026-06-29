package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	ZoneID       uint           `gorm:"not null" json:"zone_id"`
	LicensePlate string         `gorm:"size:15;not null" json:"license_plate"`
	Status       string         `gorm:"default:'active'" json:"status"` // active, completed, cancelled
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone *ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}

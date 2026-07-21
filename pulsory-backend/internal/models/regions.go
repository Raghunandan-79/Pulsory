package models

import (
	"github.com/google/uuid"
)

type Region struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Name string `gorm:"not null"`

	Ticks []WebsiteTick `gorm:"foreignKey:RegionID"`
}
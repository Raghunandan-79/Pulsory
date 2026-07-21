package models

import (
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	URL string `gorm:"not null"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	TimeAdded time.Time `gorm:"not null"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Ticks []WebsiteTick `gorm:"foreignKey:WebsiteID"`
}
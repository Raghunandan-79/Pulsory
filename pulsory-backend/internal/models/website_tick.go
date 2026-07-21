package models

import (
	"time"

	"github.com/google/uuid"
)

type WebsiteStatus string

const (
	StatusUp      WebsiteStatus = "Up"
	StatusDown    WebsiteStatus = "Down"
	StatusUnknown WebsiteStatus = "Unknown"
)

type WebsiteTick struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	ResponseTimeMS int `gorm:"not null"`

	Status WebsiteStatus `gorm:"type:varchar(20);not null"`

	RegionID uuid.UUID `gorm:"type:uuid;not null"`
	WebsiteID uuid.UUID `gorm:"type:uuid;not null"`

	Region  Region  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Website Website `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
}
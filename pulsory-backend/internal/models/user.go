package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Websites []Website `gorm:"foreignKey:UserID"`
}
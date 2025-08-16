package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Number      int            `gorm:"uniqueIndex" json:"number"`
	Type        string         `json:"type"`
	Price       float64        `json:"price"`
	IsAvailable bool           `json:"is_available"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Bookings  []Booking  `gorm:"foreignKey:RoomID"`
	Feedbacks []Feedback `gorm:"foreignKey:RoomNum;references:Number"`
}

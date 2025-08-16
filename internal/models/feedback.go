package models

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"` // Unique Feedback ID
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`                        // Who submitted
	UserName  string         `json:"user_name"`                                                // Display name
	Message   string         `json:"message"`                                                  // Feedback content
	BookingID string         `gorm:"type:uuid" json:"booking_id,omitempty"`                    // Optional, link to booking
	RoomNum   int            `json:"room_num,omitempty"`                                       // Optional, linked room
	Rating    int            `json:"rating,omitempty"`                                         // Optional rating (1-5)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID"`
	Booking Booking `gorm:"foreignKey:BookingID"`
}

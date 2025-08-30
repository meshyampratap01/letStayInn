package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	BookingStatusBooked    = "Booked"
	BookingStatusCancelled = "Cancelled"
	BookingStatusCompleted = "Completed"
)

type Booking struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string         `gorm:"type:uuid;not null;constraint:OnDelete:SET NULL;" json:"user_id"`
	RoomID     string         `gorm:"type:uuid;not null;constraint:OnDelete:SET NULL;" json:"room_id"`
	RoomNum    int            `json:"room_num"`
	CheckIn    time.Time      `json:"check_in"`
	CheckOut   time.Time      `json:"check_out"`
	Status     string         `json:"status"`
	FoodReq    bool           `json:"food_req"`
	CleanReq   bool           `json:"clean_req"`
	FeedbackID []string       `gorm:"-" json:"feedback_id"` 
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User      User       `gorm:"foreignKey:UserID"`
	Room      Room       `gorm:"foreignKey:RoomID"`
	Feedbacks []Feedback `gorm:"foreignKey:BookingID"`
}

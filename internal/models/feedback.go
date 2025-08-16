package models

import "time"

type Feedback struct {
	ID        string    `json:"id"`          // Unique Feedback ID
	UserID    string    `json:"user_id"`     // Who submitted
	UserName  string    `json:"user_name"`   // Display name
	Message   string    `json:"message"`     // Feedback content
	CreatedAt time.Time `json:"created_at"`  // When it was submitted
	BookingID string    `json:"booking_id,omitempty"` // Optional, link to booking
	RoomNum   int       `json:"room_num,omitempty"`   // Optional, linked room
	Rating    int       `json:"rating,omitempty"`     // Optional rating (1-5)
}

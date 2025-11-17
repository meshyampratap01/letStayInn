package models

import (
	"time"
)

type Feedback struct {
	PK        string    `dynamodbav:"pk" json:"pk"`
	SK        string    `dynamodbav:"sk" json:"sk"`
	ID        string    `dynamodbav:"id" json:"id"`
	UserID    string    `dynamodbav:"user_id" json:"user_id"`
	UserName  string    `dynamodbav:"user_name" json:"user_name"`
	Message   string    `dynamodbav:"message" json:"message"`
	BookingID string    `dynamodbav:"booking_id" json:"booking_id,omitempty"`
	RoomNum   int       `dynamodbav:"room_num" json:"room_num,omitempty"`
	Rating    int       `dynamodbav:"rating" json:"rating,omitempty"`
	CreatedAt time.Time `dynamodbav:"created_at" json:"created_at"`
}

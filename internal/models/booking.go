package models

import (
	"time"
)

const (
	BookingStatusBooked    = "Booked"
	BookingStatusCancelled = "Cancelled"
	BookingStatusCompleted = "Completed"
)

type Booking struct {
	PK       string    `dynamodbav:"pk" json:"pk"`
	SK       string    `dynamodbav:"sk" json:"sk"`
	ID       string    `dynamodbav:"id" json:"id"`
	UserID   string    `dynamodbav:"user_id" json:"user_id"`
	RoomID   string    `dynamodbav:"room_id" json:"room_id"`
	RoomNum  int       `dynamodbav:"room_num" json:"room_num"`
	CheckIn  time.Time `dynamodbav:"check_in" json:"check_in"`
	CheckOut time.Time `dynamodbav:"check_out" json:"check_out"`
	Status   string    `dynamodbav:"status" json:"status"`
	FoodReq  bool      `dynamodbav:"food_req" json:"food_req"`
	CleanReq bool      `dynamodbav:"clean_req" json:"clean_req"`
}

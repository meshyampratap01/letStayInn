package models

import "time"

const (
	BookingStatusBooked    = "Booked"
	BookingStatusCancelled = "Cancelled"
	BookingStatusCompleted = "Completed"
)

type Booking struct{
	ID 			string		`json:"id"`
	UserID		string		`json:"user_id"`
	RoomID		string		`json:"room_id"`
	RoomNum		int			`json:"room_num"`
	CheckIn 	time.Time	`json:"check_in"`
	CheckOut 	time.Time	`json:"check_out"`
	Status 		string		`json:"status"`
	FoodReq 	bool		`json:"food_req"`
	CleanReq 	bool		`json:"clean_req"`
	FeedbackID 	[]string	`json:"feedback_id"`
}
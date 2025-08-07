package models

import "time"

type BookingInfo struct {
	ID         string
	GuestName  string
	RoomNumber int
	CheckIn    time.Time
	CheckOut   time.Time
}

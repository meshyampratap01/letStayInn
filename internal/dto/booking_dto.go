package dto

type BookingDTO struct {
	ID         string `json:"id"`
	RoomNumber int `json:"room_number"`
	Status     string `json:"status"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
}

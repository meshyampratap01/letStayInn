package dto

type FeedbackDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	RoomNum   int `json:"room_num"`
	BookingID string `json:"booking_id"`
	Rating    int    `json:"rating"`
}

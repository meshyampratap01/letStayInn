package dto

type RoomDTO struct {
	Number      int  `json:"number"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	Description string  `json:"description"`
}

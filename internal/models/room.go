package models

type RoomType string

const (
	RoomTypeStandard  RoomType = "Standard"
	RoomTypeDeluxe    RoomType = "Deluxe"
	RoomTypeSuite     RoomType = "Suite"
	RoomTypeExecutive RoomType = "Executive"
)

type Room struct {
	ID          string    `dynamodbav:"pk" json:"pk"`
	SK          string    `dynamodbav:"sk" json:"sk"`
	Number      int       `dynamodbav:"number" json:"number"`
	Type        RoomType  `dynamodbav:"room_type" json:"room_type"`
	Price       float64   `dynamodbav:"price" json:"price"`
	IsAvailable bool      `dynamodbav:"is_available" json:"is_available"`
	Description string    `dynamodbav:"description" json:"description"`

	// Relationships
	Bookings  []Booking  `json:"-"`
	Feedbacks []Feedback `json:"-"`
}

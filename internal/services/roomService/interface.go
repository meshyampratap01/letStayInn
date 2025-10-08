package roomService

import "github.com/meshyampratap01/letStayInn/internal/models"
//go:generate mockgen -source=interface.go -destination=../../mocks/mock_roomService.go -package=mocks

type IRoomService interface {
	GetAvailableRooms() ([]models.Room, error)
	GetTotalRooms() (int, error)
	GetTotalAvailableRooms() (int,error)
	GetAllRooms() ([]models.Room, error)
	AddRoom(number int, roomType string, price float64, isAvailable bool, description string) error
	UpdateRoom(number int,roomType string, price float64, isAvailable bool, description string) error
	DeleteRoom(number int) error
	RoomExists(number int) (bool, error)
}

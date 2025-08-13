package roomService

import "github.com/meshyampratap01/letStayInn/internal/models"

type IRoomService interface {
	GetAvailableRooms() ([]models.Room, error)
	GetTotalRooms() (int, error)
	GetTotalAvailableRooms() (int,error)
	GetAllRooms() ([]models.Room, error)
	AddRoom(number int, roomType string, price float64, isAvailable bool, description string) error
	UpdateRoom(number int, choice int, roomType string, price float64, isAvailable bool, description string) error
	DeleteRoom(number int) error
}

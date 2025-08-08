package roomRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type IRoomRepository interface {
	GetAllRooms() ([]models.Room, error)
	SaveRooms(rooms []models.Room) error
	GetAvailableRooms() ([]models.Room, error)
	AddRoom(room models.Room) error
	GetRoomNumberByBookingID(string) (string, error)
}

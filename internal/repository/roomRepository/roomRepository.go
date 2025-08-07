package roomRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type RoomRepository struct{}

func NewRoomRepository() IRoomRepository {
	return &RoomRepository{}
}

func (rr *RoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := storage.ReadJson(config.RoomsFile, &rooms)
	return rooms, err
}

func (rr *RoomRepository) SaveRooms(rooms []models.Room) error {
	return storage.WriteJson(config.RoomsFile, rooms)
}

func (rr *RoomRepository) GetAvailableRooms() ([]models.Room, error) {
	rooms, err := rr.GetAllRooms()
	if err != nil {
		return nil, err
	}
	var available []models.Room
	for _, r := range rooms {
		if r.IsAvailable {
			available = append(available, r)
		}
	}
	return available, nil
}

func (r *RoomRepository) AddRoom(room models.Room) error {
	rooms, err := r.GetAllRooms()
	if err != nil {
		return err
	}
	rooms = append(rooms, room)
	return r.SaveRooms(rooms)
}

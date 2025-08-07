package roomService

import (
	"errors"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type RoomService struct {
	roomRepo roomRepository.IRoomRepository
}

func NewRoomService(roomRepo roomRepository.IRoomRepository) RoomServiceManager {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

func (rs *RoomService) GetTotalRooms() (int, error) {
	rooms, err := rs.roomRepo.GetAllRooms()
	if err != nil {
		return 0, err
	}
	return len(rooms), nil
}

func (r *RoomService) GetAvailableRooms() ([]models.Room, error) {
	return r.roomRepo.GetAvailableRooms()
}

func (r *RoomService) GetTotalAvailableRooms() (int, error) {
	totalRooms, err := r.roomRepo.GetAvailableRooms()
	return len(totalRooms), err
}

func (r *RoomService) GetAllRooms() ([]models.Room, error) {
	return r.roomRepo.GetAllRooms()
}

func (r *RoomService) AddRoom(number int, roomType string, price float64, isAvailable bool, description string) error {
	newRoom := models.Room{
		ID:          utils.NewUUID(),
		Number:      number,
		Type:        roomType,
		Price:       price,
		IsAvailable: isAvailable,
		Description: description,
	}
	return r.roomRepo.AddRoom(newRoom)
}

func (r *RoomService) UpdateRoom(number int, choice int, roomType string, price float64, isAvailable bool, description string) error {
	rooms, err := r.roomRepo.GetAllRooms()
	if err != nil {
		return err
	}

	updated := false
	for i, room := range rooms {
		if room.Number == number {
			switch choice {
			case 1:
				rooms[i].Type = roomType
			case 2:
				rooms[i].Price = price
			case 3:
				rooms[i].IsAvailable = isAvailable
			case 4:
				rooms[i].Description = description
			default:
				return errors.New("invalid update choice")
			}
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("room not found")
	}

	return r.roomRepo.SaveRooms(rooms)
}

func (r *RoomService) DeleteRoom(number int) error {
	rooms, err := r.roomRepo.GetAllRooms()
	if err != nil {
		return err
	}

	found := false
	newRooms := make([]models.Room, 0, len(rooms))

	for _, room := range rooms {
		if room.Number == number {
			found = true
			continue // skip this room (deleting it)
		}
		newRooms = append(newRooms, room)
	}

	if !found {
		return errors.New("room not found")
	}

	return r.roomRepo.SaveRooms(newRooms)
}


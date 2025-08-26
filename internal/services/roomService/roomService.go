package roomService

import (
	"errors"
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

var (
	ErrInvalidRoomNumber   = errors.New("room number must be greater than 0")
	ErrInvalidRoomType     = errors.New("invalid room type")
	ErrInvalidRoomPrice    = errors.New("room price must be greater than 0")
	ErrEmptyRoomDesc       = errors.New("room description cannot be empty")
	ErrDuplicateRoomNumber = errors.New("room number already exists")
	ErrRoomNotFound        = errors.New("room not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidChoice       = errors.New("invalid update choice")
	ErrPersistenceFailed   = errors.New("failed to save rooms")
)

type RoomService struct {
	roomRepo roomRepository.IRoomRepository
}

func NewRoomService(roomRepo roomRepository.IRoomRepository) IRoomService {
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
	// --- Validation ---
	if number <= 0 {
		return ErrInvalidRoomNumber
	}
	if !validators.IsValidRoomType(roomType) {
		return fmt.Errorf("%w: %s", ErrInvalidRoomType, roomType)
	}
	if price <= 0 {
		return ErrInvalidRoomPrice
	}
	if description == "" {
		return ErrEmptyRoomDesc
	}

	// --- Duplicate check ---
	exists, err := r.roomRepo.RoomExists(number)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateRoomNumber
	}

	// --- Create room ---
	newRoom := models.Room{
		ID:          utils.NewUUID(),
		Number:      number,
		Type:        models.RoomType(roomType),
		Price:       price,
		IsAvailable: isAvailable,
		Description: description,
	}
	return r.roomRepo.AddRoom(newRoom)
}

func (r *RoomService) UpdateRoom(number int, choice int, roomType string, price float64, isAvailable bool, description string) error {
	// Validate update choice
	if choice < 1 || choice > 4 {
		return fmt.Errorf("%w: %d (must be between 1 and 4)", ErrInvalidChoice, choice)
	}

	rooms, err := r.roomRepo.GetAllRooms()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPersistenceFailed, err)
	}

	updated := false
	for i, room := range rooms {
		if room.Number == number {
			switch choice {
			case 1: // Update Type
				if roomType == "" {
					return fmt.Errorf("%w: room type cannot be empty", ErrInvalidInput)
				}
				rooms[i].Type = models.RoomType(roomType)

			case 2: // Update Price
				if price <= 0 {
					return fmt.Errorf("%w: price %.2f must be greater than 0", ErrInvalidInput, price)
				}
				rooms[i].Price = price

			case 3: // Update Availability
				rooms[i].IsAvailable = isAvailable

			case 4: // Update Description
				if len(description) < 5 {
					return fmt.Errorf("%w: description must be at least 5 characters long", ErrInvalidInput)
				}
				rooms[i].Description = description
			}
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("%w: room number %d", ErrRoomNotFound, number)
	}

	// Persist changes
	if err := r.roomRepo.SaveRooms(rooms); err != nil {
		return fmt.Errorf("%w: %v", ErrPersistenceFailed, err)
	}

	return nil
}

func (r *RoomService) DeleteRoom(number int) error {
	return r.roomRepo.DeleteRoomByNumber(number)
}

func (s *RoomService) RoomExists(number int) (bool, error) {
	return s.roomRepo.RoomExists(number)
}

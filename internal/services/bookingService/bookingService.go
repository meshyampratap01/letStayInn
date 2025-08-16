package bookingService

import (
	"context"
	"errors"
	"fmt"
	"time"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type BookingService struct {
	bookingRepo bookingRepository.BookingRepository
	roomRepo    roomRepository.IRoomRepository
	userRepo	userRepository.UserRepository
}

func NewBookingService(bookingRepo bookingRepository.BookingRepository, roomRepo roomRepository.IRoomRepository,userRepo userRepository.UserRepository) IBookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		userRepo: 	 userRepo,
	}
}

func (s *BookingService) BookRoom(ctx context.Context, roomNum int, checkInStr, checkOutStr string) error {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return err
	}

	var selected *models.Room
	for i := range rooms {
		if rooms[i].Number == roomNum && rooms[i].IsAvailable {
			selected = &rooms[i]
			break
		}
	}
	if selected == nil {
		return errors.New("room not available")
	}

	checkIn, err := time.Parse("02-01-2006", checkInStr)
	if err != nil {
		return errors.New("invalid check-in date")
	}
	checkOut, err := time.Parse("02-01-2006", checkOutStr)
	if err != nil {
		return errors.New("invalid check-out date")
	}

	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return err
	}

	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("invalid or missing user ID in context")
	}

	newBooking := models.Booking{
		ID:       utils.NewUUID(),
		UserID:   userID,
		RoomID:   selected.ID,
		RoomNum:  selected.Number,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Status:   models.BookingStatusBooked,
	}
	bookings = append(bookings, newBooking)

	if err := s.bookingRepo.SaveBookings(bookings); err != nil {
		return err
	}


	for i := range rooms {
		if rooms[i].ID == selected.ID {
			rooms[i].IsAvailable = false
			break
		}
	}
	if err := s.roomRepo.SaveRooms(rooms); err != nil {
		return err
	}

	return nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID string) error {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return err
	}
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("invalid or missing user ID in context")
	}

	var updated bool
	for i := range bookings {
		b := &bookings[i]
		if b.ID == bookingID && b.UserID == userID && b.Status == models.BookingStatusBooked {
			b.Status = models.BookingStatusCancelled
			updated = true

			rooms, err := s.roomRepo.GetAllRooms()
			if err != nil {
				return err
			}
			for j := range rooms {
				if rooms[j].ID == b.RoomID {
					rooms[j].IsAvailable = true
					break
				}
			}
			if err := s.roomRepo.SaveRooms(rooms); err != nil {
				return err
			}
			break
		}
	}
	if !updated {
		return errors.New("booking not found or already cancelled")
	}

	return s.bookingRepo.SaveBookings(bookings)
}

func (s *BookingService) GetActiveBookings() ([]models.Booking, error) {
	return s.bookingRepo.GetActiveBookings()
}


func (s *BookingService) GetUserActiveBookings(ctx context.Context) ([]models.Booking, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return nil,fmt.Errorf("invalid or missing user ID in context")
	}
	bookings, err := s.bookingRepo.GetBookingsByUserID(userID)
	if err != nil {
		return nil, err
	}

	active := []models.Booking{}
	for _, b := range bookings {
		if b.Status == models.BookingStatusBooked {
			active = append(active, b)
		}
	}
	return active, nil
}


func (s *BookingService) GetBookingIDByRoomNumber(roomNumber int) (string, error) {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return "", err
	}

	for _, b := range bookings {
		if b.RoomNum == roomNumber && b.Status == models.BookingStatusBooked {
			return b.ID, nil
		}
	}
	return "", nil
}


func (bs *BookingService) IsRoomBooked(roomNumber int) (bool, error) {
	return bs.bookingRepo.CheckRoomBooked(roomNumber)
}

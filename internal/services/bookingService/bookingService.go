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
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type BookingService struct {
	bookingRepo bookingRepository.BookingRepository
	roomRepo    roomRepository.IRoomRepository
	userRepo    userRepository.UserRepository
	serviceRepo	serviceRequestRepository.ServiceRequestRepository
}

func NewBookingService(bookingRepo bookingRepository.BookingRepository, roomRepo roomRepository.IRoomRepository, userRepo userRepository.UserRepository,serviceRepo	serviceRequestRepository.ServiceRequestRepository) IBookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *BookingService) BookRoom(ctx context.Context, roomNum int, checkInStr, checkOutStr string) error {
	if roomNum <= 0 {
		return errors.New("invalid room number")
	}

	validCheckIn, err := validators.ValidateDate(checkInStr)
	if err != nil {
		return fmt.Errorf("invalid check-in date: %w", err)
	}

	_, err = validators.ValidateCheckoutDate(validCheckIn, checkOutStr)
	if err != nil {
		return fmt.Errorf("invalid check-out date: %w", err)
	}

	room, err := s.roomRepo.GetRoomByNumber(roomNum)
	if err != nil {
		return err
	}
	if room == nil || !room.IsAvailable {
		return errors.New("room not available")
	}

	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("invalid or missing user ID in context")
	}

	layout := "02-01-2006"
	checkInTime, _ := time.Parse(layout, checkInStr)
	checkOutTime, _ := time.Parse(layout, checkOutStr)

	newBooking := models.Booking{
		ID:       utils.NewUUID(),
		UserID:   userID,
		RoomID:   room.ID,
		RoomNum:  room.Number,
		CheckIn:  checkInTime,
		CheckOut: checkOutTime,
		Status:   models.BookingStatusBooked,
	}

	if err := s.bookingRepo.SaveBooking(newBooking); err != nil {
		return err
	}

	room.IsAvailable = false
	if err := s.roomRepo.SaveRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID string) error {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("invalid or missing user ID in context")
	}

	booking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return errors.New("booking not found")
	}

	if booking.UserID != userID || booking.Status != models.BookingStatusBooked {
		return errors.New("booking not found or already cancelled")
	}

	err = s.serviceRepo.DeleteRoomRequests(booking.RoomNum)
	if err!= nil{
		return err
	}

	booking.Status = models.BookingStatusCancelled
	if err := s.bookingRepo.UpdateBooking(*booking); err != nil {
		return err
	}

	room, err := s.roomRepo.GetRoomByNumber(booking.RoomNum)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	room.IsAvailable = true
	if err := s.roomRepo.SaveRoom(room); err != nil {
		return err
	}


	return nil
}

func (s *BookingService) GetRoomByNumber(roomNum int) (*models.Room, error) {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return nil, err
	}
	for i := range rooms {
		if rooms[i].Number == roomNum {
			return &rooms[i], nil
		}
	}
	return nil, fmt.Errorf("room not found")
}

func (s *BookingService) GetActiveBookings() ([]models.Booking, error) {
	return s.bookingRepo.GetActiveBookings()
}

func (s *BookingService) GetUserActiveBookings(ctx context.Context) ([]models.Booking, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing user ID in context")
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
	// Check if room is booked
	isBooked, err := s.bookingRepo.CheckRoomBooked(roomNumber)
	if err != nil {
		return "", err
	}
	if !isBooked {
		return "", nil
	}

	// Room is booked, fetch active bookings to find the booking ID
	// This is a fallback when we need the booking ID; in production consider adding a direct query
	active, err := s.bookingRepo.GetActiveBookings()
	if err != nil {
		return "", err
	}

	for _, b := range active {
		if b.RoomNum == roomNumber && b.Status == models.BookingStatusBooked {
			return b.ID, nil
		}
	}
	return "", nil
}

func (bs *BookingService) IsRoomBooked(roomNumber int) (bool, error) {
	return bs.bookingRepo.CheckRoomBooked(roomNumber)
}

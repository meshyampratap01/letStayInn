package bookingRepository

import (
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileBookingRepository struct{}

func NewFileBookingRepository() BookingRepository {
	return &FileBookingRepository{}
}

func (db *FileBookingRepository) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	err := storage.ReadJson(config.BookingsFile, &bookings)
	return bookings, err
}

func (db *FileBookingRepository) SaveBookings(bookings []models.Booking) error {
	return storage.WriteJson(config.BookingsFile, bookings)
}

func (db *FileBookingRepository) GetBookingsByUserID(userID string) ([]models.Booking, error) {
	bookings, err := db.GetAllBookings()
	if err != nil {
		return nil, err
	}
	var result []models.Booking
	for _, b := range bookings {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (db *FileBookingRepository) UpdateBooking(updated models.Booking) error {
	bookings, err := db.GetAllBookings()
	if err != nil {
		return err
	}

	for i := range bookings {
		if bookings[i].ID == updated.ID {
			bookings[i] = updated
			break
		}
	}
	return db.SaveBookings(bookings)
}

func (br *FileBookingRepository) GetBookingByID(bookingID string) (*models.Booking, error) {
	bookings, err := br.GetAllBookings()
	if err != nil {
		return nil, fmt.Errorf("failed to load bookings: %w", err)
	}

	for _, booking := range bookings {
		if booking.ID == bookingID {
			return &booking, nil
		}
	}

	return nil, fmt.Errorf("booking with ID %s not found", bookingID)
}

func (r *FileBookingRepository) GetActiveBookings() ([]models.Booking, error) {
	bookings, err := r.GetAllBookings()
	if err != nil {
		return nil, err
	}

	activeBookings := []models.Booking{}
	for _, b := range bookings {
		if b.Status != models.BookingStatusCancelled {
			activeBookings = append(activeBookings, b)
		}
	}

	return activeBookings, nil
}

func (br *FileBookingRepository) CheckRoomBooked(roomNumber int) (bool, error) {
	bookings, err := br.GetAllBookings()
	if err != nil {
		return false, err
	}

	for _, b := range bookings {
		if b.RoomNum == roomNumber && b.Status == models.BookingStatusBooked {
			return true, nil
		}
	}

	return false, nil
}


package bookingRepository

import (

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileBookingRepository struct{}

func NewFileBookingRepository() BookingRepository{
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



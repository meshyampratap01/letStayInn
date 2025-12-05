//go:generate mockgen -source=interface.go -destination=../../mocks/mock_bookingRepository.go -package=mocks
package bookingRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type BookingRepository interface {
	GetAllBookings() ([]models.Booking, error)
	GetActiveBookings() ([]models.Booking, error)
	SaveBooking(booking models.Booking) error
	GetBookingsByUserID(userID string) ([]models.Booking, error)
	UpdateBooking(updated models.Booking) error
	GetBookingByID(string) (*models.Booking, error)
	CheckRoomBooked(roomNumber int) (bool, error)
	GetExpiredBookings() ([]models.Booking, error) 
}

package bookingRepository


import "github.com/meshyampratap01/letStayInn/internal/models"

type BookingRepository interface {
	GetAllBookings() ([]models.Booking, error)
	SaveBookings(bookings []models.Booking) error
	GetBookingsByUserID(userID string) ([]models.Booking, error)
	UpdateBooking(updated models.Booking) error
}
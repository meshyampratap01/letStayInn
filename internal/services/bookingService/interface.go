package bookingService

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

// BookingManager defines the business logic interface for booking-related operations.
type BookingManager interface {
	BookRoom(ctx context.Context, roomNum int, checkInStr, checkOutStr string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetUserActiveBookings(ctx context.Context) ([]models.Booking, error)
	GetAllBookingsWithGuests() ([]models.BookingInfo, error)
}



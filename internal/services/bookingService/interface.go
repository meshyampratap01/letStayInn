package bookingService

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

type IBookingService interface {
	BookRoom(ctx context.Context, roomNum int, checkInStr, checkOutStr string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetUserActiveBookings(ctx context.Context) ([]models.Booking, error)
	GetAllBookingsWithGuests() ([]models.BookingInfo, error)
	GetBookingIDByRoomNumber(roomNumber int) (string, error)
}



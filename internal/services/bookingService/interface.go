package bookingService

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/mock_bookingService.go -package=mocks

type IBookingService interface {
	GetRoomByNumber(roomNum int) (*models.Room, error)
	BookRoom(ctx context.Context, roomNum int, checkInStr, checkOutStr string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetUserActiveBookings(ctx context.Context) ([]models.Booking, error)
	GetActiveBookings() ([]models.Booking, error)
	GetBookingIDByRoomNumber(roomNumber int) (string, error)
	IsRoomBooked(roomNumber int) (bool, error)
	UpdateCompletedBookings() error
}

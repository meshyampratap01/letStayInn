package bookingService

import (
	"context"
	"errors"
	"testing"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"go.uber.org/mock/gomock"
)

func TestBookRoom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

	rooms := []models.Room{
		{ID: "room-1", Number: 101, IsAvailable: true},
	}
	bookings := []models.Booking{}

	// Expectations
	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil).Times(1)
	mockBookingRepo.EXPECT().SaveBookings(gomock.Any()).Return(nil).Times(1)
	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil).Times(1)

	err := service.BookRoom(ctx, 101, "25-08-2027", "27-08-2027")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBookRoom_InvalidRoomNumber(t *testing.T) {
    service := NewBookingService(nil, nil, nil)
    ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

    err := service.BookRoom(ctx, 0, "25-08-2027", "27-08-2027")
    if err == nil || err.Error() != "invalid room number" {
        t.Errorf("expected invalid room number error, got %v", err)
    }
}

func TestBookRoom_MissingUserID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
    mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
    mockUserRepo := mocks.NewMockUserRepository(ctrl)

    service := NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

    ctx := context.Background() // no user id
    rooms := []models.Room{{ID: "room-1", Number: 101, IsAvailable: true}}
    mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)
    mockBookingRepo.EXPECT().GetAllBookings().Return([]models.Booking{}, nil).Times(1)

    err := service.BookRoom(ctx, 101, "25-08-2027", "27-08-2027")
    if err == nil || err.Error() != "invalid or missing user ID in context" {
        t.Errorf("expected invalid/missing user ID error, got %v", err)
    }
}



func TestBookRoom_RoomNotAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

	rooms := []models.Room{
		{ID: "room-1", Number: 101, IsAvailable: false}, // Not available
	}

	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)

	err := service.BookRoom(ctx, 101, "25-08-2027", "27-08-2027")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "room not available" {
		t.Errorf("expected 'room not available', got %v", err)
	}
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

	bookings := []models.Booking{
		{
			ID:     "booking-1",
			UserID: "user-123",
			RoomID: "room-1",
			Status: models.BookingStatusBooked,
		},
	}
	rooms := []models.Room{
		{ID: "room-1", Number: 101, IsAvailable: false},
	}

	// Expectations
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil).Times(1)
	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)
	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil).Times(1)
	mockBookingRepo.EXPECT().SaveBookings(gomock.Any()).Return(nil).Times(1)

	err := service.CancelBooking(ctx, "booking-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCancelBooking_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

	bookings := []models.Booking{} // no bookings
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil).Times(1)

	err := service.CancelBooking(ctx, "non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "booking not found or already cancelled" {
		t.Errorf("expected 'booking not found or already cancelled', got %v", err)
	}
}

func TestCancelBooking_MissingUserID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
    mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
    mockUserRepo := mocks.NewMockUserRepository(ctrl)

    service := NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

    ctx := context.Background() // no user id
    mockBookingRepo.EXPECT().GetAllBookings().Return([]models.Booking{}, nil).Times(1)

    err := service.CancelBooking(ctx, "booking-1")
    if err == nil || err.Error() != "invalid or missing user ID in context" {
        t.Errorf("expected invalid/missing user ID error, got %v", err)
    }
}

func TestGetRoomByNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	rooms := []models.Room{{ID: "room-1", Number: 101, IsAvailable: true}}
	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)

	room, err := service.GetRoomByNumber(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if room.Number != 101 {
		t.Errorf("expected room number 101, got %d", room.Number)
	}
}

func TestGetRoomByNumber_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	rooms := []models.Room{{ID: "room-1", Number: 102, IsAvailable: true}}
	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil).Times(1)

	room, err := service.GetRoomByNumber(101)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if room != nil {
		t.Errorf("expected nil room, got %+v", room)
	}
}

func TestGetRoomByNumber_Error(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
    service := NewBookingService(nil, mockRoomRepo, nil)

    mockRoomRepo.EXPECT().GetAllRooms().Return(nil, errors.New("db error")).Times(1)

    _, err := service.GetRoomByNumber(101)
    if err == nil || err.Error() != "db error" {
        t.Errorf("expected db error, got %v", err)
    }
}

func TestGetActiveBookings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	expected := []models.Booking{{ID: "b1", Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetActiveBookings().Return(expected, nil).Times(1)

	bookings, err := service.GetActiveBookings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bookings) != 1 || bookings[0].ID != "b1" {
		t.Errorf("unexpected bookings: %+v", bookings)
	}
}

func TestGetActiveBookings_Error(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
    service := NewBookingService(mockBookingRepo, nil, nil)

    mockBookingRepo.EXPECT().GetActiveBookings().Return(nil, errors.New("db error")).Times(1)

    _, err := service.GetActiveBookings()
    if err == nil || err.Error() != "db error" {
        t.Errorf("expected db error, got %v", err)
    }
}

func TestGetUserActiveBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")
	bookings := []models.Booking{
		{ID: "b1", UserID: "user-123", Status: models.BookingStatusBooked},
		{ID: "b2", UserID: "user-123", Status: models.BookingStatusCancelled},
	}
	mockBookingRepo.EXPECT().GetBookingsByUserID("user-123").Return(bookings, nil).Times(1)

	active, err := service.GetUserActiveBookings(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(active) != 1 || active[0].ID != "b1" {
		t.Errorf("unexpected active bookings: %+v", active)
	}
}

func TestGetUserActiveBookings_MissingUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)
	ctx := context.Background() // no userID

	active, err := service.GetUserActiveBookings(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if active != nil {
		t.Errorf("expected nil active bookings, got %+v", active)
	}
}

func TestGetUserActiveBookings_Error(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
    service := NewBookingService(mockBookingRepo, nil, nil)

    ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

    mockBookingRepo.EXPECT().GetBookingsByUserID("user-123").Return(nil, errors.New("db error")).Times(1)

    _, err := service.GetUserActiveBookings(ctx)
    if err == nil || err.Error() != "db error" {
        t.Errorf("expected db error, got %v", err)
    }
}

func TestGetBookingIDByRoomNumber_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	bookings := []models.Booking{
		{ID: "b1", RoomNum: 101, Status: models.BookingStatusBooked},
		{ID: "b2", RoomNum: 102, Status: models.BookingStatusCancelled},
	}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil).Times(1)

	id, err := service.GetBookingIDByRoomNumber(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "b1" {
		t.Errorf("expected booking id b1, got %s", id)
	}
}

func TestGetBookingIDByRoomNumber_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)

	bookings := []models.Booking{
		{ID: "b2", RoomNum: 102, Status: models.BookingStatusCancelled},
	}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil).Times(1)

	id, err := service.GetBookingIDByRoomNumber(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "" {
		t.Errorf("expected empty id, got %s", id)
	}
}


func TestGetBookingIDByRoomNumber_Error(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
    service := NewBookingService(mockBookingRepo, nil, nil)

    mockBookingRepo.EXPECT().GetAllBookings().Return(nil, errors.New("db error")).Times(1)

    _, err := service.GetBookingIDByRoomNumber(101)
    if err == nil || err.Error() != "db error" {
        t.Errorf("expected db error, got %v", err)
    }
}

func TestIsRoomBooked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)
	mockBookingRepo.EXPECT().CheckRoomBooked(101).Return(true, nil).Times(1)

	isBooked, err := service.IsRoomBooked(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isBooked {
		t.Errorf("expected true, got false")
	}
}

func TestIsRoomBooked_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service :=  NewBookingService(mockBookingRepo, mockRoomRepo, mockUserRepo)
	mockBookingRepo.EXPECT().CheckRoomBooked(101).Return(false, errors.New("db error")).Times(1)

	isBooked, err := service.IsRoomBooked(101)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if isBooked {
		t.Errorf("expected false, got true")
	}
}

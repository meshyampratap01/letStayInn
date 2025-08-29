package bookingRepository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() { sqlDB.Close() }
	return gdb, mock, cleanup
}

func TestGetAllBookings(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "room_num", "status"}).
		AddRow("1", "u1", 101, models.BookingStatusBooked).
		AddRow("2", "u2", 102, models.BookingStatusBooked)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings"`)).
		WillReturnRows(rows)

	bookings, err := repo.GetAllBookings()
	if err != nil || len(bookings) != 2 {
		t.Errorf("expected 2 bookings, got %v, err=%v", bookings, err)
	}
}

func TestSaveBookings(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	bookings := []models.Booking{
		{ID: "1", UserID: "u1", RoomNum: 101, Status: models.BookingStatusBooked},
		{ID: "2", UserID: "u2", RoomNum: 102, Status: models.BookingStatusBooked},
	}

	for range bookings {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	if err := repo.SaveBookings(bookings); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetBookingsByUserID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "room_num", "status"}).
		AddRow("1", "u1", 101, models.BookingStatusBooked)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE user_id = $1`)).
		WithArgs("u1").
		WillReturnRows(rows)

	bookings, err := repo.GetBookingsByUserID("u1")
	if err != nil || len(bookings) != 1 {
		t.Errorf("expected 1 booking, got %v, err=%v", bookings, err)
	}
}

func TestUpdateBooking(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	booking := models.Booking{ID: "1", UserID: "u1", RoomNum: 101, Status: models.BookingStatusBooked}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.UpdateBooking(booking); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetBookingByID_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "room_num", "status"}).
		AddRow("1", "u1", 101, models.BookingStatusBooked)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE id = $1 AND "bookings"."deleted_at" IS NULL ORDER BY "bookings"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	booking, err := repo.GetBookingByID("1")
	if err != nil || booking.ID != "1" {
		t.Errorf("expected booking with id 1, got %v, err=%v", booking, err)
	}
}

func TestGetBookingByID_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE id = $1 AND "bookings"."deleted_at" IS NULL ORDER BY "bookings"."id" LIMIT $2`)).
		WithArgs("404", 1).
		WillReturnError(errors.New("not found"))

	booking, err := repo.GetBookingByID("404")
	if err == nil || booking != nil {
		t.Errorf("expected error, got booking=%v, err=%v", booking, err)
	}
}

func TestGetActiveBookings(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "room_num", "status"}).
		AddRow("1", "u1", 101, models.BookingStatusBooked).
		AddRow("2", "u2", 102, models.BookingStatusBooked)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE status != $1 AND "bookings"."deleted_at" IS NULL`)).
		WithArgs(models.BookingStatusCancelled).
		WillReturnRows(rows)

	bookings, err := repo.GetActiveBookings()
	if err != nil || len(bookings) != 2 {
		t.Errorf("expected 2 bookings, got %v, err=%v", bookings, err)
	}
}

func TestCheckRoomBooked(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormBookingRepository(db)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bookings" WHERE (room_num = $1 AND status = $2) AND "bookings"."deleted_at" IS NULL`)).
		WithArgs(101, models.BookingStatusBooked).
		WillReturnRows(rows)

	booked, err := repo.CheckRoomBooked(101)
	if err != nil || !booked {
		t.Errorf("expected booked=true, got %v, err=%v", booked, err)
	}
}

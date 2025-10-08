package roomRepository

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	dialector := postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestGetAllRooms(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	rows := sqlmock.NewRows([]string{"id", "number", "type", "price", "is_available"}).
		AddRow("1", 101, "Deluxe", 100.0, true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	rooms, err := repo.GetAllRooms()
	if err != nil || len(rooms) != 1 {
		t.Errorf("expected 1 room, got %v, err %v", rooms, err)
	}

	// error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("query failed"))
	_, err = repo.GetAllRooms()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestSaveRooms(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	rooms := []models.Room{
		{ID: "1", Number: 101, Type: "Deluxe", Price: 100},
		{ID: "2", Number: 102, Type: "Suite", Price: 200},
	}

	// first room
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET`)).
		WithArgs(
			rooms[0].Number, rooms[0].Type, rooms[0].Price, rooms[0].IsAvailable,
			rooms[0].Description, sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), rooms[0].ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// second room
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET`)).
		WithArgs(
			rooms[1].Number, rooms[1].Type, rooms[1].Price, rooms[1].IsAvailable,
			rooms[1].Description, sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), rooms[1].ID,
		).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	if err := repo.SaveRooms(rooms); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error on second save
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET`)).
		WillReturnError(errors.New("update failed"))
	mock.ExpectRollback()

	err := repo.SaveRooms([]models.Room{{ID: "3", Number: 103}})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// func TestDeleteRoomByNumber(t *testing.T) {
// 	db, mock, cleanup := setupMockDB(t)
// 	defer cleanup()
// 	repo := NewGormRoomRepository(db)

// 	// success: fetch then delete
// 	roomID := "room-uuid"
// 	rows := sqlmock.NewRows([]string{"id", "number"}).AddRow(roomID, 101)
// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE number = $1 AND "rooms"."deleted_at" IS NULL ORDER BY "rooms"."id" LIMIT $2`)).
// 		WithArgs(101, 1).
// 		WillReturnRows(rows)
// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET "deleted_at"`)).
// 		WithArgs(sqlmock.AnyArg(), roomID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	err := repo.DeleteRoomByNumber(101)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}

// 	// failure: room not found
// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE number = $1 AND "rooms"."deleted_at" IS NULL ORDER BY "rooms"."id" LIMIT $2`)).
// 		WithArgs(999, 1).
// 		WillReturnError(errors.New("not found"))
// 	err = repo.DeleteRoomByNumber(999)
// 	if err == nil {
// 		t.Errorf("expected error, got nil")
// 	}
// }

func TestGetAvailableRooms(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	rows := sqlmock.NewRows([]string{"id", "number", "type", "price", "is_available"}).
		AddRow("1", 101, "Deluxe", 100.0, true)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "rooms" WHERE is_available = $1 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(true).
		WillReturnRows(rows)

	rooms, err := repo.GetAvailableRooms()
	if err != nil || len(rooms) != 1 {
		t.Errorf("expected available rooms, got %v, err %v", rooms, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "rooms" WHERE is_available = $1 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(true).
		WillReturnError(errors.New("query failed"))

	_, err = repo.GetAvailableRooms()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestAddRoom(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	room := models.Room{ID: "1", Number: 101, Type: "Deluxe"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WithArgs(
			room.Number, room.Type, room.Price, room.IsAvailable,
			room.Description, sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), room.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	if err := repo.AddRoom(room); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	if err := repo.AddRoom(room); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetRoomNumberByBookingID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	rows := sqlmock.NewRows([]string{"id", "room_num"}).
		AddRow("b1", 101)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookings" WHERE id = $1 AND "bookings"."deleted_at" IS NULL ORDER BY "bookings"."id" LIMIT $2`)).
		WithArgs("b1", 1).
		WillReturnRows(rows)

	num, err := repo.GetRoomNumberByBookingID("b1")
	if err != nil || num != "101" {
		t.Errorf("expected room number 101, got %s, err %v", num, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookings" WHERE id = $1 AND "bookings"."deleted_at" IS NULL ORDER BY "bookings"."id" LIMIT $2`)).
		WithArgs("b1", 1).
		WillReturnError(errors.New("not found"))

	_, err = repo.GetRoomNumberByBookingID("b1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRoomExists(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormRoomRepository(db)

	// exists
	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "rooms" WHERE number = $1 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(101).
		WillReturnRows(rows)

	exists, err := repo.RoomExists(101)
	if err != nil || !exists {
		t.Errorf("expected true, got %v, err %v", exists, err)
	}

	// not exists
	rows = sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "rooms" WHERE number = $1 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(101).
		WillReturnRows(rows)

	exists, err = repo.RoomExists(101)
	if err != nil || exists {
		t.Errorf("expected false, got %v, err %v", exists, err)
	}

	// error case
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "rooms" WHERE number = $1 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(101).
		WillReturnError(errors.New("count failed"))

	_, err = repo.RoomExists(101)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

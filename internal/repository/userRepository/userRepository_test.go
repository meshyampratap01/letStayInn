package userRepository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockUserDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestGetAllUsers(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "available"}).
		AddRow("1", "Alice", "alice@test.com", true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers()
	if err != nil || len(users) != 1 {
		t.Errorf("expected 1 user, got %v, err %v", users, err)
	}

	// error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("query failed"))
	_, err = repo.GetAllUsers()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestFindUserByEmail(t *testing.T) {
	repo := NewGormUserRepository(nil)

	users := []models.User{
		{ID: "1", Name: "Alice", Email: "alice@test.com"},
		{ID: "2", Name: "Bob", Email: "bob@test.com"},
	}

	found := repo.FindUserByEmail(users, "alice@test.com")
	if found == nil || found.Name != "Alice" {
		t.Errorf("expected Alice, got %+v", found)
	}

	notFound := repo.FindUserByEmail(users, "charlie@test.com")
	if notFound != nil {
		t.Errorf("expected nil, got %+v", notFound)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow("1", "Alice", "alice@test.com")
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("alice@test.com", 1).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("alice@test.com")
	if err != nil || user.Name != "Alice" {
		t.Errorf("expected Alice, got %+v, err %v", user, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("missing@test.com", 1).
		WillReturnError(errors.New("not found"))

	_, err = repo.GetUserByEmail("missing@test.com")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestSaveUser(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	user := models.User{ID: "1", Name: "Alice", Email: "alice@test.com"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(
			user.Name,
			user.Email,
			sqlmock.AnyArg(), // password
			sqlmock.AnyArg(), // role
			user.Available,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			user.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	if err := repo.SaveUser(user); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error case
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	if err := repo.SaveUser(user); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	user := &models.User{ID: "1", Name: "Alice", Email: "alice@test.com"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs(
			user.Name,
			user.Email,
			sqlmock.AnyArg(), // password
			sqlmock.AnyArg(), // role
			user.Available,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			user.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.UpdateUser(user); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSaveAllUsers(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	users := []models.User{
		{ID: "1", Name: "Alice", Email: "alice@test.com"},
		{ID: "2", Name: "Bob", Email: "bob@test.com"},
	}

	// success
	for _, u := range users {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(
				u.Name,
				u.Email,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				u.Available,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				u.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	if err := repo.SaveAllUsers(users); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// failure on one user
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WillReturnError(errors.New("update failed"))
	mock.ExpectRollback()

	err := repo.SaveAllUsers([]models.User{{ID: "3", Name: "Charlie"}})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("1", "Alice")
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID("1")
	if err != nil || user.Name != "Alice" {
		t.Errorf("expected Alice, got %+v, err %v", user, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("2", 1).
		WillReturnError(errors.New("not found"))

	_, err = repo.GetUserByID("2")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestToggleStaffAvailability(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	// success
	rows := sqlmock.NewRows([]string{"id", "available"}).
		AddRow("1", false)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true,
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.ToggleStaffAvailability("1"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// not found
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("999", 1).
		WillReturnError(errors.New("not found"))

	if err := repo.ToggleStaffAvailability("999"); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetStaffAvailability(t *testing.T) {
	db, mock, cleanup := setupMockUserDB(t)
	defer cleanup()
	repo := NewGormUserRepository(db)

	// available
	rows := sqlmock.NewRows([]string{"id", "available"}).
		AddRow("1", true)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	avail, err := repo.GetStaffAvailability("1")
	if err != nil || !avail {
		t.Errorf("expected true, got %v, err %v", avail, err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("999", 1).
		WillReturnError(errors.New("not found"))

	_, err = repo.GetStaffAvailability("999")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// func TestDeleteUserByID(t *testing.T) {
// 	db, mock, cleanup := setupMockUserDB(t)
// 	defer cleanup()
// 	repo := NewGormUserRepository(db)

// 	// success: fetch then delete
// 	userID := "user-uuid"
// 	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(userID, "Alice")
// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
// 		WithArgs(userID, 1).
// 		WillReturnRows(rows)
// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"`)).
// 		WithArgs(sqlmock.AnyArg(), userID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	err := repo.DeleteUserByID(userID)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}

// 	// failure: user not found
// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
// 		WithArgs("missing-id", 1).
// 		WillReturnError(errors.New("not found"))
// 	err = repo.DeleteUserByID("missing-id")
// 	if err == nil {
// 		t.Errorf("expected error, got nil")
// 	}
// }

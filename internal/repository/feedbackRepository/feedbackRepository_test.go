package feedbackRepository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// helper to set up sqlmock + gorm
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

func TestSaveFeedback_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormFeedbackRepository(db)

	feedback := models.Feedback{
		ID:      "1",
		UserID:  "u1",
		Message: "Great stay!",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "feedbacks"`)).
		WithArgs(
			feedback.UserID,
			feedback.UserName,
			feedback.Message,
			feedback.BookingID,
			feedback.RoomNum,
			feedback.Rating,
			sqlmock.AnyArg(), 
			sqlmock.AnyArg(), 
			sqlmock.AnyArg(), 
			feedback.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	err := repo.SaveFeedback(feedback)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
func TestSaveFeedback_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormFeedbackRepository(db)

	feedback := models.Feedback{
		ID:      "1",
		UserID:  "u1",
		Message: "Great stay!",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "feedbacks"`)).
		WithArgs(
			feedback.UserID,
			feedback.UserName,
			feedback.Message,
			feedback.BookingID,
			feedback.RoomNum,
			feedback.Rating,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			feedback.ID,
		).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err := repo.SaveFeedback(feedback)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetAllFeedback_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormFeedbackRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "message"}).
		AddRow("1", "u1", "Great stay!").
		AddRow("2", "u2", "Nice service!")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "feedbacks"`)).
		WillReturnRows(rows)

	feedbacks, err := repo.GetAllFeedback()
	if err != nil || len(feedbacks) != 2 {
		t.Errorf("expected 2 feedbacks, got %v, err=%v", feedbacks, err)
	}
}

func TestGetAllFeedback_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewGormFeedbackRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "feedbacks"`)).
		WillReturnError(errors.New("query failed"))

	feedbacks, err := repo.GetAllFeedback()
	if err == nil || feedbacks != nil {
		t.Errorf("expected error, got feedbacks=%v err=%v", feedbacks, err)
	}
}

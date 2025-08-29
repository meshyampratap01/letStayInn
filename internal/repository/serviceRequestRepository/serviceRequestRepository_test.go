package serviceRequestRepository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *GormServiceRequestRepository) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening gorm db: %v", err)
	}
	repo := NewGormServiceRequestRepository(gdb).(*GormServiceRequestRepository)
	return gdb, mock, repo
}

func TestLoadServiceRequests(t *testing.T) {
	_, mock, repo := setup(t)

	rows := sqlmock.NewRows([]string{"id", "room_num"}).
		AddRow("1", 101).
		AddRow("2", 102)

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*deleted_at.*`).
		WillReturnRows(rows)

	reqs, err := repo.LoadServiceRequests()
	if err != nil || len(reqs) != 2 {
		t.Errorf("expected 2 requests, got %v err=%v", len(reqs), err)
	}
}

func TestLoadServiceRequests_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*deleted_at.*`).
		WillReturnError(errors.New("db error"))

	_, err := repo.LoadServiceRequests()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSaveServiceRequest(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "service_requests"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
	mock.ExpectCommit()

	err := repo.SaveServiceRequest(models.ServiceRequest{ID: "123"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSaveServiceRequest_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "service_requests"`).
		WillReturnError(errors.New("insert fail"))
	mock.ExpectRollback()

	err := repo.SaveServiceRequest(models.ServiceRequest{ID: "123"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetUnassignedRequests(t *testing.T) {
	_, mock, repo := setup(t)

	rows := sqlmock.NewRows([]string{"id", "is_assigned", "status"}).
		AddRow("1", false, "pending").
		AddRow("2", false, "assigned")

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*deleted_at.*`).
		WillReturnRows(rows)

	reqs, err := repo.GetUnassignedRequests()
	if err != nil || len(reqs) != 2 {
		t.Errorf("expected 2 reqs, got %d, err=%v", len(reqs), err)
	}
}

func TestGetServiceRequestByRoomNum(t *testing.T) {
	_, mock, repo := setup(t)

	rows := sqlmock.NewRows([]string{"id", "room_num"}).
		AddRow("123", 101)

	// GORM sends args: room_num, limit(1)
	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*room_num.*deleted_at.*`).
		WithArgs(101, 1).
		WillReturnRows(rows)

	req, err := repo.GetServiceRequestByRoomNum(101)
	if err != nil || req.RoomNum != 101 {
		t.Errorf("expected room_num 101, got %v, err=%v", req, err)
	}
}

func TestGetServiceRequestByRoomNum_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*room_num.*deleted_at.*`).
		WithArgs(101, 1).
		WillReturnError(errors.New("not found"))

	_, err := repo.GetServiceRequestByRoomNum(101)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetServiceRequestByReqID(t *testing.T) {
	_, mock, repo := setup(t)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("123")
	// GORM sends args: id, limit(1)
	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*id.*deleted_at.*`).
		WithArgs("123", 1).
		WillReturnRows(rows)

	req, err := repo.GetServiceRequestByReqID("123")
	if err != nil || req.ID != "123" {
		t.Errorf("expected id=123, got %v, err=%v", req, err)
	}
}

func TestGetServiceRequestByReqID_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*id.*deleted_at.*`).
		WithArgs("404", 1).
		WillReturnError(errors.New("not found"))

	_, err := repo.GetServiceRequestByReqID("404")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateServiceRequest(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "service_requests"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := &models.ServiceRequest{ID: "321"}
	err := repo.UpdateServiceRequest(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateServiceRequest_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "service_requests"`).
		WillReturnError(errors.New("update fail"))
	mock.ExpectRollback()

	req := &models.ServiceRequest{ID: "321"}
	err := repo.UpdateServiceRequest(req)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAssignedServiceRequests(t *testing.T) {
	_, mock, repo := setup(t)

	rows := sqlmock.NewRows([]string{"id", "assigned_to"}).
		AddRow("1", "emp1").
		AddRow("2", "emp1")

	mock.ExpectQuery(`SELECT .* FROM "service_requests" WHERE .*assigned_to.*deleted_at.*`).
		WithArgs("emp1").
		WillReturnRows(rows)

	reqs, err := repo.GetAssignedServiceRequests("emp1")
	if err != nil || len(reqs) != 2 {
		t.Errorf("expected 2 reqs, got %d, err=%v", len(reqs), err)
	}
}

func TestUpdateIsAssigned(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	// GORM sends args: is_assigned, updated_at, id
	mock.ExpectExec(`UPDATE "service_requests"`).
		WithArgs(true, sqlmock.AnyArg(), "req123").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateIsAssigned("req123", true)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
}

func TestUpdateIsAssigned_Error(t *testing.T) {
	_, mock, repo := setup(t)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "service_requests"`).
		WithArgs(false, sqlmock.AnyArg(), "req123").
		WillReturnError(errors.New("fail"))
	mock.ExpectRollback()

	err := repo.UpdateIsAssigned("req123", false)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

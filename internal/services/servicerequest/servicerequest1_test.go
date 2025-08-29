package servicerequest

import (
	"context"
	"errors"
	"testing"
	"time"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"go.uber.org/mock/gomock"
)

func setupServiceRequestService(ctrl *gomock.Controller) (*ServiceRequestService, *mocks.MockBookingRepository, *mocks.MockServiceRequestRepository) {
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	service := NewServiceRequestService(mockBookingRepo, mockServiceRequestRepo)
	return service.(*ServiceRequestService), mockBookingRepo, mockServiceRequestRepo
}

func TestServiceRequestGetter_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(&bookings[0], nil)
	mockBookingRepo.EXPECT().UpdateBooking(gomock.Any()).Return(nil)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return([]models.ServiceRequest{}, nil)
	mockServiceRequestRepo.EXPECT().SaveServiceRequest(gomock.Any()).Return(nil)
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServiceRequestGetter_NoBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, _ := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	mockBookingRepo.EXPECT().GetAllBookings().Return([]models.Booking{}, nil)
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "you don't have any active or completed booking for room 101" {
		t.Errorf("expected no booking error, got %v", err)
	}
}

func TestServiceRequestGetter_GetAllBookingsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, _ := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	mockBookingRepo.EXPECT().GetAllBookings().Return(nil, errors.New("load error"))
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "failed to load bookings: load error" {
		t.Errorf("expected load error, got %v", err)
	}
}

func TestServiceRequestGetter_GetBookingByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, _ := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(nil, errors.New("not found"))
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "error fetching the required booking: not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestServiceRequestGetter_UpdateBookingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, _ := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(&bookings[0], nil)
	mockBookingRepo.EXPECT().UpdateBooking(gomock.Any()).Return(errors.New("update error"))
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "failed to update booking request flags: update error" {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestServiceRequestGetter_LoadServiceRequestsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(&bookings[0], nil)
	mockBookingRepo.EXPECT().UpdateBooking(gomock.Any()).Return(nil)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return(nil, errors.New("load error"))
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "failed to load service requests: load error" {
		t.Errorf("expected load error, got %v", err)
	}
}

func TestServiceRequestGetter_PreviousRequestInProgress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(&bookings[0], nil)
	mockBookingRepo.EXPECT().UpdateBooking(gomock.Any()).Return(nil)
	latest := models.ServiceRequest{UserID: "user-1", RoomNum: 101, Type: models.ServiceTypeFood, Status: models.ServiceStatusPending, CreatedAt: time.Now()}
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return([]models.ServiceRequest{latest}, nil)
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "your previous Food request for this room is still in progress" {
		t.Errorf("expected previous request in progress error, got %v", err)
	}
}

func TestServiceRequestGetter_SaveServiceRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockBookingRepo, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-1")
	bookings := []models.Booking{{ID: "b1", UserID: "user-1", RoomNum: 101, Status: models.BookingStatusBooked}}
	mockBookingRepo.EXPECT().GetAllBookings().Return(bookings, nil)
	mockBookingRepo.EXPECT().GetBookingByID("b1").Return(&bookings[0], nil)
	mockBookingRepo.EXPECT().UpdateBooking(gomock.Any()).Return(nil)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return([]models.ServiceRequest{}, nil)
	mockServiceRequestRepo.EXPECT().SaveServiceRequest(gomock.Any()).Return(errors.New("save error"))
	err := service.ServiceRequestGetter(ctx, 101, models.ServiceTypeFood, "details")
	if err == nil || err.Error() != "failed to save request: save error" {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestGetPendingRequestCount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	pending := []models.ServiceRequest{{Status: models.ServiceStatusPending}, {Status: models.ServiceStatusDone}}
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return(pending, nil)
	count, err := service.GetPendingRequestCount()
	if err != nil || count != 1 {
		t.Errorf("expected 1 pending, got %v, err %v", count, err)
	}
}

func TestGetPendingRequestCount_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return(nil, errors.New("load error"))
	_, err := service.GetPendingRequestCount()
	if err == nil || err.Error() != "load error" {
		t.Errorf("expected load error, got %v", err)
	}
}

func TestGetUnassignedServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().GetUnassignedRequests().Return([]models.ServiceRequest{{ID: "r1"}}, nil)
	reqs, err := service.GetUnassignedServiceRequest()
	if err != nil || len(reqs) != 1 {
		t.Errorf("expected 1 unassigned, got %v, err %v", len(reqs), err)
	}
}

func TestGetUnassignedServiceRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().GetUnassignedRequests().Return(nil, errors.New("unassigned error"))
	_, err := service.GetUnassignedServiceRequest()
	if err == nil || err.Error() != "unassigned error" {
		t.Errorf("expected unassigned error, got %v", err)
	}
}

func TestCancelServiceRequestByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	req := &models.ServiceRequest{ID: "r1"}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(req, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(nil)
	if err := service.CancelServiceRequestByID("r1"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCancelServiceRequestByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(nil, errors.New("not found"))
	if err := service.CancelServiceRequestByID("r1"); err == nil || err.Error() != "service request not found: not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestCancelServiceRequestByID_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	req := &models.ServiceRequest{ID: "r1"}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(req, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(errors.New("update error"))
	if err := service.CancelServiceRequestByID("r1"); err == nil || err.Error() != "update error" {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestUpdateServiceRequestStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	req := &models.ServiceRequest{ID: "r1"}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(req, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(nil)
	if err := service.UpdateServiceRequestStatus("r1", models.ServiceStatusDone); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateServiceRequestStatus_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(nil, errors.New("not found"))
	if err := service.UpdateServiceRequestStatus("r1", models.ServiceStatusDone); err == nil || err.Error() != "service request not found: not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestUpdateServiceRequestStatus_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	req := &models.ServiceRequest{ID: "r1"}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("r1").Return(req, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(errors.New("update error"))
	if err := service.UpdateServiceRequestStatus("r1", models.ServiceStatusDone); err == nil || err.Error() != "update error" {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestUpdateServiceRequestAssignment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().UpdateIsAssigned("r1", true).Return(nil)
	if err := service.UpdateServiceRequestAssignment("r1", true); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateServiceRequestAssignment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo := setupServiceRequestService(ctrl)
	mockServiceRequestRepo.EXPECT().UpdateIsAssigned("r1", true).Return(errors.New("assign error"))
	if err := service.UpdateServiceRequestAssignment("r1", true); err == nil || err.Error() != "assign error" {
		t.Errorf("expected assign error, got %v", err)
	}
}

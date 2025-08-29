package employeeService

import (
	"errors"
	"reflect"
	"testing"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/mocks"

	gomock "go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*gomock.Controller,
	*mocks.MockUserRepository,
	*mocks.MockIRoomRepository,
	*mocks.MockBookingRepository,
	*mocks.MockServiceRequestRepository,
	IEmployeeService,
	) {

	ctrl := gomock.NewController(t)

	userRepo := mocks.NewMockUserRepository(ctrl)
	roomRepo := mocks.NewMockIRoomRepository(ctrl)
	bookingRepo := mocks.NewMockBookingRepository(ctrl)
	serviceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	svc :=  NewEmployeeService(userRepo, roomRepo, bookingRepo, serviceRequestRepo)

	return ctrl, userRepo, roomRepo, bookingRepo, serviceRequestRepo,svc
}

func TestGetAssignedServiceRequests(t *testing.T) {
	ctrl, _, _, _, serviceReqRepo, svc := setup(t)
	defer ctrl.Finish()

	expected := []models.ServiceRequest{{ID: "1", Type: models.ServiceTypeCleaning}}
	serviceReqRepo.EXPECT().GetAssignedServiceRequests("emp1").Return(expected, nil)

	got, err := svc.GetAssignedServiceRequests("emp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestUpdateServiceRequestStatus(t *testing.T) {
	tests := []struct {
		name        string
		requests    []models.ServiceRequest
		requestID   string
		newStatus   models.ServiceStatus
		mockSetup   func(userRepo *mocks.MockUserRepository, bookingRepo *mocks.MockBookingRepository)
		expectError bool
	}{
		{
			name:        "request not found",
			requests:    []models.ServiceRequest{},
			requestID:   "invalid",
			newStatus:   models.ServiceStatusInProgress,
			mockSetup:   func(u *mocks.MockUserRepository, b *mocks.MockBookingRepository) {},
			expectError: true,
		},
		{
			name: "mark as done updates user availability and booking",
			requests: []models.ServiceRequest{{
				ID:         "req1",
				Status:     models.ServiceStatusInProgress,
				AssignedTo: "emp1",
				BookingID:  "b1",
				Type:       models.ServiceTypeFood,
			}},
			requestID: "req1",
			newStatus: models.ServiceStatusDone,
			mockSetup: func(u *mocks.MockUserRepository, b *mocks.MockBookingRepository) {
				u.EXPECT().GetUserByID("emp1").Return(&models.User{ID: "emp1"}, nil)
				u.EXPECT().UpdateUser(gomock.Any()).Return(nil)

				b.EXPECT().GetBookingByID("b1").Return(&models.Booking{ID: "b1", FoodReq: true}, nil)
				b.EXPECT().UpdateBooking(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, userRepo, _, bookingRepo, serviceReqRepo, svc := setup(t)
			defer ctrl.Finish()

			serviceReqRepo.EXPECT().LoadServiceRequests().Return(tt.requests, nil)
			if len(tt.requests) > 0 {
				serviceReqRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(nil)
			}

			tt.mockSetup(userRepo, bookingRepo)

			err := svc.UpdateServiceRequestStatus(tt.requestID, tt.newStatus)
			if tt.expectError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestToggleAvailability(t *testing.T) {
	ctrl, userRepo, _, _, _, svc := setup(t)
	defer ctrl.Finish()

	userRepo.EXPECT().ToggleStaffAvailability("u1").Return(nil)

	if err := svc.ToggleAvailability("u1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetAvailability(t *testing.T) {
	ctrl, userRepo, _, _, _, svc := setup(t)
	defer ctrl.Finish()

	userRepo.EXPECT().GetStaffAvailability("u1").Return(true, nil)

	avail, err := svc.GetAvailability("u1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !avail {
		t.Errorf("expected true, got false")
	}
}

func TestGetRoomNumberByBookingID(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(roomRepo *mocks.MockIRoomRepository, bookingRepo *mocks.MockBookingRepository)
		expectError bool
	}{
		{
			name: "booking not found",
			mockSetup: func(r *mocks.MockIRoomRepository, b *mocks.MockBookingRepository) {
				b.EXPECT().GetBookingByID("b1").Return(nil, errors.New("not found"))
			},
			expectError: true,
		},
		{
			name: "room not found",
			mockSetup: func(r *mocks.MockIRoomRepository, b *mocks.MockBookingRepository) {
				b.EXPECT().GetBookingByID("b1").Return(&models.Booking{ID: "b1"}, nil)
				r.EXPECT().GetRoomNumberByBookingID("b1").Return("", errors.New("room not found"))
			},
			expectError: true,
		},
		{
			name: "success",
			mockSetup: func(r *mocks.MockIRoomRepository, b *mocks.MockBookingRepository) {
				b.EXPECT().GetBookingByID("b1").Return(&models.Booking{ID: "b1"}, nil)
				r.EXPECT().GetRoomNumberByBookingID("b1").Return("101", nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, _, roomRepo, bookingRepo, _, svc := setup(t)
			defer ctrl.Finish()

			tt.mockSetup(roomRepo, bookingRepo)

			_, err := svc.GetRoomNumberByBookingID("b1")
			if tt.expectError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

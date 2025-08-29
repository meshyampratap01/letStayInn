package managerservice

import (
	"errors"
	"testing"

	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"go.uber.org/mock/gomock"
)

func setupManagerService(ctrl *gomock.Controller) (*ManagerService, *mocks.MockUserRepository, *mocks.MockServiceRequestRepository, *mocks.MockIRoomRepository, *mocks.MockBookingRepository, *mocks.MockFeedbackRepository) {
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockFeedbackRepo := mocks.NewMockFeedbackRepository(ctrl)
	service := NewManagerService(mockUserRepo, mockServiceRequestRepo, mockRoomRepo, mockBookingRepo, mockFeedbackRepo)
	return service.(*ManagerService), mockUserRepo, mockServiceRequestRepo, mockRoomRepo, mockBookingRepo, mockFeedbackRepo
}

func TestUpdateEmployeeAvailability_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Email: "a@b.com", Role: models.RoleKitchenStaff, Available: false}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	mockUserRepo.EXPECT().SaveAllUsers([]models.User{{Email: "a@b.com", Role: models.RoleKitchenStaff, Available: true}}).Return(nil)
	if err := service.UpdateEmployeeAvailability("a@b.com", true); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateEmployeeAvailability_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)
	if err := service.UpdateEmployeeAvailability("x@y.com", true); err == nil || err.Error() != "employee with email x@y.com not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestGetAllEmployees_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Role: models.RoleKitchenStaff}, {Role: models.RoleCleaningStaff}, {Role: models.RoleManager}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	emps, err := service.GetAllEmployees()
	if err != nil || len(emps) != 2 {
		t.Errorf("expected 2 employees, got %v, err %v", len(emps), err)
	}
}

func TestGetTotalEmployees_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Role: models.RoleKitchenStaff}, {Role: models.RoleCleaningStaff}, {Role: models.RoleManager}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	count, err := service.GetTotalEmployees()
	if err != nil || count != 2 {
		t.Errorf("expected 2, got %v, err %v", count, err)
	}
}

func TestDeleteEmployeeByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Email: "a@b.com", Role: models.RoleKitchenStaff}, {Email: "b@c.com", Role: models.RoleManager}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	mockUserRepo.EXPECT().SaveAllUsers([]models.User{{Email: "b@c.com", Role: models.RoleManager}}).Return(nil)
	if err := service.DeleteEmployeeByEmail("a@b.com"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteEmployeeByEmail_NotEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Email: "a@b.com", Role: models.RoleManager}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	if err := service.DeleteEmployeeByEmail("a@b.com"); err == nil || err.Error() != "user with this email is not an employee" {
		t.Errorf("expected not employee error, got %v", err)
	}
}

func TestDeleteEmployeeByEmail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)
	if err := service.DeleteEmployeeByEmail("x@y.com"); err == nil || err.Error() != "employee not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestGetAvailableStaffByTaskType_Cleaning(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Role: models.RoleCleaningStaff, Available: true}, {Role: models.RoleCleaningStaff, Available: false}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	staff, err := service.GetAvailableStaffByTaskType(string(models.ServiceTypeCleaning))
	if err != nil || len(staff) != 1 {
		t.Errorf("expected 1 available cleaning staff, got %v, err %v", len(staff), err)
	}
}

func TestGetAvailableStaffByTaskType_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, _, _, _ := setupManagerService(ctrl)
	_, err := service.GetAvailableStaffByTaskType("invalid")
	if err == nil || err.Error() != "invalid task type" {
		t.Errorf("expected invalid task type error, got %v", err)
	}
}

func TestAssignServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, mockServiceRequestRepo, _, _, _ := setupManagerService(ctrl)
	req := &models.ServiceRequest{ID: "req1"}
	emp := &models.User{ID: "emp1", Available: true}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("req1").Return(req, nil)
	mockUserRepo.EXPECT().GetUserByID("emp1").Return(emp, nil)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(nil)
	if err := service.AssignServiceRequest("req1", "emp1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetHotelReport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, mockServiceRequestRepo, mockRoomRepo, mockBookingRepo, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}}, nil)
	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{}}, nil)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{{}}, nil)
	mockBookingRepo.EXPECT().GetAllBookings().Return([]models.Booking{{}}, nil)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return([]models.ServiceRequest{{IsAssigned: false, Status: models.ServiceStatusPending}}, nil)
	report, err := service.GetHotelReport()
	if err != nil || report.TotalRooms != 1 || report.AvailableRooms != 1 || report.TotalStaff != 0 || report.TotalBookings != 1 || report.UnassignedRequests != 1 {
		t.Errorf("unexpected report values or error: %v, %+v", err, report)
	}
}

// Additional error path tests for coverage
func TestUpdateEmployeeAvailability_GetAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	if err := service.UpdateEmployeeAvailability("a@b.com", true); err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestUpdateEmployeeAvailability_SaveAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Email: "a@b.com", Role: models.RoleKitchenStaff, Available: false}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	mockUserRepo.EXPECT().SaveAllUsers(gomock.Any()).Return(errors.New("save error"))
	if err := service.UpdateEmployeeAvailability("a@b.com", true); err == nil || err.Error() != "save error" {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestGetAllEmployees_GetAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	_, err := service.GetAllEmployees()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetTotalEmployees_GetAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	_, err := service.GetTotalEmployees()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestDeleteEmployeeByEmail_GetAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	if err := service.DeleteEmployeeByEmail("a@b.com"); err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestDeleteEmployeeByEmail_SaveAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	users := []models.User{{Email: "a@b.com", Role: models.RoleKitchenStaff}, {Email: "b@c.com", Role: models.RoleManager}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	mockUserRepo.EXPECT().SaveAllUsers(gomock.Any()).Return(errors.New("save error"))
	if err := service.DeleteEmployeeByEmail("a@b.com"); err == nil || err.Error() != "save error" {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestGetAvailableStaffByTaskType_GetAllUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, _, _, _ := setupManagerService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	_, err := service.GetAvailableStaffByTaskType(string(models.ServiceTypeCleaning))
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestAssignServiceRequest_GetServiceRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, mockServiceRequestRepo, _, _, _ := setupManagerService(ctrl)
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("req1").Return(nil, errors.New("not found"))
	if err := service.AssignServiceRequest("req1", "emp1"); err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestAssignServiceRequest_UpdateServiceRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, mockServiceRequestRepo, _, _, _ := setupManagerService(ctrl)
	req := &models.ServiceRequest{ID: "req1"}
	emp := &models.User{ID: "emp1", Available: true}
	mockServiceRequestRepo.EXPECT().GetServiceRequestByReqID("req1").Return(req, nil)
	mockUserRepo.EXPECT().GetUserByID("emp1").Return(emp, nil)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(gomock.Any()).Return(errors.New("update error"))
	if err := service.AssignServiceRequest("req1", "emp1"); err == nil || err.Error() != "update error" {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestGetHotelReport_GetAllRoomsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, mockRoomRepo, _, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return(nil, errors.New("rooms error"))
	_, err := service.GetHotelReport()
	if err == nil || err.Error() != "error fetching rooms: rooms error" {
		t.Errorf("expected rooms error, got %v", err)
	}
}

func TestGetHotelReport_GetAvailableRoomsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, mockRoomRepo, _, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}}, nil)
	mockRoomRepo.EXPECT().GetAvailableRooms().Return(nil, errors.New("avail error"))
	_, err := service.GetHotelReport()
	if err == nil || err.Error() != "error fetching available rooms: avail error" {
		t.Errorf("expected avail error, got %v", err)
	}
}

func TestGetHotelReport_GetAllEmployeesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, mockRoomRepo, _, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}}, nil)
	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{}}, nil)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("emp error"))
	_, err := service.GetHotelReport()
	if err == nil || err.Error() != "error fetching employees: emp error" {
		t.Errorf("expected emp error, got %v", err)
	}
}

func TestGetHotelReport_GetAllBookingsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, _, mockRoomRepo, mockBookingRepo, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}}, nil)
	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{}}, nil)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{{}}, nil)
	mockBookingRepo.EXPECT().GetAllBookings().Return(nil, errors.New("bookings error"))
	_, err := service.GetHotelReport()
	if err == nil || err.Error() != "error fetching bookings: bookings error" {
		t.Errorf("expected bookings error, got %v", err)
	}
}

func TestGetHotelReport_LoadServiceRequestsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo, mockServiceRequestRepo, mockRoomRepo, mockBookingRepo, _ := setupManagerService(ctrl)
	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}}, nil)
	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{}}, nil)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{{}}, nil)
	mockBookingRepo.EXPECT().GetAllBookings().Return([]models.Booking{{}}, nil)
	mockServiceRequestRepo.EXPECT().LoadServiceRequests().Return(nil, errors.New("sr error"))
	_, err := service.GetHotelReport()
	if err == nil || err.Error() != "error fetching service requests: sr error" {
		t.Errorf("expected sr error, got %v", err)
	}
}

func TestViewAllFeedback_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, _, _, mockFeedbackRepo := setupManagerService(ctrl)
	mockFeedbackRepo.EXPECT().GetAllFeedback().Return(nil, errors.New("feedback error"))
	_, err := service.ViewAllFeedback()
	if err == nil || err.Error() != "feedback error" {
		t.Errorf("expected feedback error, got %v", err)
	}
}

func TestViewAllFeedback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _, _, _, _, mockFeedbackRepo := setupManagerService(ctrl)
	mockFeedbackRepo.EXPECT().GetAllFeedback().Return([]models.Feedback{{ID: "f1"}}, nil)
	feedbacks, err := service.ViewAllFeedback()
	if err != nil || len(feedbacks) != 1 {
		t.Errorf("expected 1 feedback, got %v, err %v", len(feedbacks), err)
	}
}

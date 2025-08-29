package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// setup for manager handler
func setupManager(t *testing.T) (*gomock.Controller,
	*mocks.MockIRoomService,
	*mocks.MockIBookingService,
	*mocks.MockIUserService,
	*mocks.MockIServiceRequestService,
	*mocks.MockIManagerService,
	*ManagerHandler) {

	ctrl := gomock.NewController(t)

	mockRoom := mocks.NewMockIRoomService(ctrl)
	mockBooking := mocks.NewMockIBookingService(ctrl)
	mockUser := mocks.NewMockIUserService(ctrl)
	mockSR := mocks.NewMockIServiceRequestService(ctrl)
	mockManager := mocks.NewMockIManagerService(ctrl)

	h := NewManagerHandler(mockRoom, mockBooking, mockUser, mockSR, mockManager)

	// disable logging
	logger.Log = zap.NewNop()

	return ctrl, mockRoom, mockBooking, mockUser, mockSR, mockManager, h
}

// utility to make a manager-context request
func managerReq(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager"))
}

func nonManagerReq(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Guest"))
}

func TestUpdateRoomHTTP_Success(t *testing.T) {
	ctrl, mockRoom, _, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockRoom.EXPECT().UpdateRoom(101, 1, "Deluxe", 200.0, true, "Nice").Return(nil)

	body := `{"choice":1,"type":"Deluxe","price":200,"is_available":true,"description":"Nice"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/101", strings.NewReader(body))
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()

	h.UpdateRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestUpdateRoomHTTP_InvalidRole(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/101", nil)
	rec := httptest.NewRecorder()
	h.UpdateRoomHTTP(rec, nonManagerReq(req))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestDeleteRoomHTTP_RoomBooked(t *testing.T) {
	ctrl, _, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().IsRoomBooked(101).Return(true, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestDeleteRoomHTTP_Unauthorized(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()

	// not wrapped with managerReq → should fail RequireManager
	h.DeleteRoomHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestDeleteRoomHTTP_InvalidRoomNumber(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	// non-numeric
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/abc", nil)
	req.SetPathValue("roomNum", "abc")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}

	// zero or negative
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/0", nil)
	req.SetPathValue("roomNum", "0")
	rec = httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestDeleteRoomHTTP_IsRoomBookedError(t *testing.T) {
	ctrl, _, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().IsRoomBooked(101).Return(false, errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "db error") {
		t.Errorf("expected db error message, got %s", rec.Body.String())
	}
}

func TestDeleteRoomHTTP_DeleteRoomError(t *testing.T) {
	ctrl, mockRoom, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().IsRoomBooked(101).Return(false, nil)
	mockRoom.EXPECT().DeleteRoom(101).Return(errors.New("delete failed"))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "delete failed") {
		t.Errorf("expected delete failed message, got %s", rec.Body.String())
	}
}

func TestDeleteRoomHTTP_Success(t *testing.T) {
	ctrl, mockRoom, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().IsRoomBooked(101).Return(false, nil)
	mockRoom.EXPECT().DeleteRoom(101).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Room deleted successfully" {
		t.Errorf("expected success message, got %s", resp["message"])
	}
}


func TestListEmployeesHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().GetAllEmployees().Return([]models.User{
		{ID: "e1", Name: "John", Email: "j@x.com", Role: models.RoleCleaningStaff, Available: true},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	rec := httptest.NewRecorder()
	h.ListEmployeesHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestCreateEmployeeHTTP_InvalidRole(t *testing.T) {
	_, _, _, mockUser, _, _, h := setupManager(t)

	body := `{"name":"John","email":"j@x.com","password":"pass","role":"Invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}

	// also cover success
	mockUser.EXPECT().CreateEmployee("John", "j@x.com", "pass", models.RoleKitchenStaff, true).
		Return(models.User{ID: "e1"}, nil)

	body = `{"name":"John","email":"j@x.com","password":"pass","role":"KitchenStaff","available":true}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/employees", strings.NewReader(body))
	rec = httptest.NewRecorder()
	h.CreateEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rec.Code)
	}
}

func TestDeleteEmployeeHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().DeleteEmployeeByEmail("e1").Return(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/e1", nil)
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()
	h.DeleteEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestUpdateEmployeeAvailabilityHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().UpdateEmployeeAvailability("e1", true).Return(nil)
	body := `{"available":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader(body))
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()
	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestListAllBookingsHTTP_AllTrue(t *testing.T) {
	ctrl, _, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().GetActiveBookings().Return([]models.Booking{{ID: "b1"}}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings?all=true", nil)
	rec := httptest.NewRecorder()
	h.ListAllBookingsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestListUnassignedServiceRequestsHTTP(t *testing.T) {
	ctrl, _, _, _, mockSR, _, h := setupManager(t)
	defer ctrl.Finish()

	mockSR.EXPECT().GetUnassignedServiceRequest().Return([]models.ServiceRequest{{ID: "sr1"}}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-requests?unassigned=true", nil)
	rec := httptest.NewRecorder()
	h.ListUnassignedServiceRequestsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestAssignServiceRequestHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().AssignServiceRequest("sr1", "e1").Return(nil)
	body := `{"employee_id":"e1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/sr1/assign", strings.NewReader(body))
	req.SetPathValue("requestId", "sr1")
	rec := httptest.NewRecorder()
	h.AssignServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_InvalidStatus(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)
	body := `{"status":"Invalid"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/s1/status", strings.NewReader(body))
	req.SetPathValue("requestId", "s1")
	rec := httptest.NewRecorder()
	h.UpdateServiceRequestStatusHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_MissingRequestID(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests//status", strings.NewReader(body))
	// intentionally not setting requestId
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Missing service request id") {
		t.Errorf("expected missing id error, got %s", rec.Body.String())
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_InvalidBody(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	body := `{"status":` // malformed JSON
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/s1/status", strings.NewReader(body))
	req.SetPathValue("requestId", "s1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Invalid request body") {
		t.Errorf("expected invalid body error, got %s", rec.Body.String())
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_ServiceError(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	// mock service failure
	h.serviceRequestService.(*mocks.MockIServiceRequestService).
		EXPECT().
		UpdateServiceRequestStatus("s1", models.ServiceStatusPending).
		Return(errors.New("update failed"))

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/s1/status", strings.NewReader(body))
	req.SetPathValue("requestId", "s1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "update failed") {
		t.Errorf("expected update failed message, got %s", rec.Body.String())
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_Success(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	// mock service success
	h.serviceRequestService.(*mocks.MockIServiceRequestService).
		EXPECT().
		UpdateServiceRequestStatus("s1", models.ServiceStatusPending).
		Return(nil)

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/s1/status", strings.NewReader(body))
	req.SetPathValue("requestId", "s1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, managerReq(req))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Service request status updated" {
		t.Errorf("expected success message, got %s", resp["message"])
	}
}

func TestManager_UpdateServiceRequestStatusHTTP_Unauthorized(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/s1/status", strings.NewReader(body))
	req.SetPathValue("requestId", "s1")
	rec := httptest.NewRecorder()

	// don’t wrap with managerReq → should fail RequireManager
	h.UpdateServiceRequestStatusHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestCancelServiceRequestHTTP_Success(t *testing.T) {
	ctrl, _, _, _, mockSR, _, h := setupManager(t)
	defer ctrl.Finish()

	mockSR.EXPECT().CancelServiceRequestByID("sr1").Return(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/sr1", nil)
	req.SetPathValue("requestId", "sr1")
	rec := httptest.NewRecorder()
	h.CancelServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestListAllFeedbackHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().ViewAllFeedback().Return([]models.Feedback{{ID: "f1"}}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?all=true", nil)
	rec := httptest.NewRecorder()
	h.ListAllFeedbackHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestGenerateReportHTTP_Success(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupManager(t)
	defer ctrl.Finish()

	mockManager.EXPECT().GetHotelReport().Return(&models.HotelReport{AvailableRooms: 10}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	rec := httptest.NewRecorder()
	h.GenerateReportHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestAddRoomHTTP_Success(t *testing.T) {
	ctrl, mockRoom, _, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	mockRoom.EXPECT().AddRoom(201, "Deluxe", 500.0, true, "Sea view").Return(nil)
	body := `{"number":201,"type":"Deluxe","price":500,"description":"Sea view"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.AddRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rec.Code)
	}
}

func TestUpdateRoomHTTP_Errors(t *testing.T) {
	_, mockRoom, _, _, _, _, h := setupManager(t)

	// invalid body
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/101", strings.NewReader("bad json"))
	req.SetPathValue("roomNum", "101")
	rec := httptest.NewRecorder()
	h.UpdateRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// invalid room number
	body := `{"choice":1}`
	req = httptest.NewRequest(http.MethodPut, "/api/v1/rooms/abc", strings.NewReader(body))
	req.SetPathValue("roomNum", "abc")
	rec = httptest.NewRecorder()
	h.UpdateRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// update error
	mockRoom.EXPECT().UpdateRoom(101, 1, "", 0.0, false, "").Return(errors.New("fail"))
	req = httptest.NewRequest(http.MethodPut, "/api/v1/rooms/101", strings.NewReader(body))
	req.SetPathValue("roomNum", "101")
	rec = httptest.NewRecorder()
	h.UpdateRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestDeleteRoomHTTP_Errors(t *testing.T) {
	ctrl, mockRoom, mockBooking, _, _, _, h := setupManager(t)
	defer ctrl.Finish()

	// invalid room number
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/xx", nil)
	req.SetPathValue("roomNum", "xx")
	rec := httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// IsRoomBooked error
	mockBooking.EXPECT().IsRoomBooked(101).Return(false, errors.New("boom"))
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/101", nil)
	req.SetPathValue("roomNum", "101")
	rec = httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}

	// DeleteRoom error
	mockBooking.EXPECT().IsRoomBooked(102).Return(false, nil)
	mockRoom.EXPECT().DeleteRoom(102).Return(errors.New("fail"))
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/102", nil)
	req.SetPathValue("roomNum", "102")
	rec = httptest.NewRecorder()
	h.DeleteRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestListEmployeesHTTP_Error(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)
	mockManager.EXPECT().GetAllEmployees().Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	rec := httptest.NewRecorder()
	h.ListEmployeesHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}
}

func TestCreateEmployeeHTTP_Errors(t *testing.T) {
	_, _, _, mockUser, _, _, h := setupManager(t)

	// invalid body
	req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", strings.NewReader("bad"))
	rec := httptest.NewRecorder()
	h.CreateEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// CreateEmployee error
	mockUser.EXPECT().CreateEmployee("John", "j@x.com", "pass", models.RoleKitchenStaff, true).
		Return(models.User{}, errors.New("fail"))
	body := `{"name":"John","email":"j@x.com","password":"pass","role":"KitchenStaff","available":true}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/employees", strings.NewReader(body))
	rec = httptest.NewRecorder()
	h.CreateEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestDeleteEmployeeHTTP_Errors(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	// missing employeeId
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/", nil)
	rec := httptest.NewRecorder()
	h.DeleteEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// delete error
	mockManager.EXPECT().DeleteEmployeeByEmail("e1").Return(errors.New("fail"))
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/employees/e1", nil)
	req.SetPathValue("employeeId", "e1")
	rec = httptest.NewRecorder()
	h.DeleteEmployeeHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestUpdateEmployeeAvailabilityHTTP_Errors(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	// invalid body
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader("bad"))
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()
	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// update error
	mockManager.EXPECT().UpdateEmployeeAvailability("e1", false).Return(errors.New("fail"))
	body := `{"available":false}`
	req = httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader(body))
	req.SetPathValue("employeeId", "e1")
	rec = httptest.NewRecorder()
	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestUpdateEmployeeAvailabilityHTTP_MissingEmployeeID(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	body := `{"available":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees//availability", strings.NewReader(body))
	// intentionally not setting employeeId
	rec := httptest.NewRecorder()

	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Missing employee id") {
		t.Errorf("expected missing employee id error, got %s", rec.Body.String())
	}
}

func TestUpdateEmployeeAvailabilityHTTP_Unauthorized(t *testing.T) {
	_, _, _, _, _, _, h := setupManager(t)

	body := `{"available":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader(body))
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()

	// don’t wrap request with managerReq → should fail RequireManager
	h.UpdateEmployeeAvailabilityHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}

func TestUpdateEmployeeAvailabilityHTTP_SuccessTrue(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	mockManager.EXPECT().UpdateEmployeeAvailability("e1", true).Return(nil)

	body := `{"available":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader(body))
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()

	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Employee availability updated" {
		t.Errorf("expected success message, got %s", resp["message"])
	}
}

func TestUpdateEmployeeAvailabilityHTTP_SuccessFalse(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	mockManager.EXPECT().UpdateEmployeeAvailability("e1", false).Return(nil)

	body := `{"available":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/e1/availability", strings.NewReader(body))
	req.SetPathValue("employeeId", "e1")
	rec := httptest.NewRecorder()

	h.UpdateEmployeeAvailabilityHTTP(rec, managerReq(req))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Employee availability updated" {
		t.Errorf("expected success message, got %s", resp["message"])
	}
}


func TestListAllBookingsHTTP_ErrorAndElse(t *testing.T) {
	_, _, mockBooking, _, _, _, h := setupManager(t)

	// error branch
	mockBooking.EXPECT().GetActiveBookings().Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings?all=true", nil)
	rec := httptest.NewRecorder()
	h.ListAllBookingsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}

	// all=false branch
	req = httptest.NewRequest(http.MethodGet, "/api/v1/bookings?all=false", nil)
	rec = httptest.NewRecorder()
	h.ListAllBookingsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rec.Code)
	}
}

func TestListUnassignedServiceRequestsHTTP_Branches(t *testing.T) {
	_, _, _, _, mockSR, _, h := setupManager(t)

	// error branch
	mockSR.EXPECT().GetUnassignedServiceRequest().Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-requests?unassigned=true", nil)
	rec := httptest.NewRecorder()
	h.ListUnassignedServiceRequestsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}

	// unassigned=false branch
	req = httptest.NewRequest(http.MethodGet, "/api/v1/service-requests?unassigned=false", nil)
	rec = httptest.NewRecorder()
	h.ListUnassignedServiceRequestsHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rec.Code)
	}
}

func TestAssignServiceRequestHTTP_Errors(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	// missing requestId
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests//assign", nil)
	rec := httptest.NewRecorder()
	h.AssignServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// invalid body
	req = httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/sr1/assign", strings.NewReader("bad"))
	req.SetPathValue("requestId", "sr1")
	rec = httptest.NewRecorder()
	h.AssignServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// service error
	mockManager.EXPECT().AssignServiceRequest("sr1", "e1").Return(errors.New("fail"))
	body := `{"employee_id":"e1"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/sr1/assign", strings.NewReader(body))
	req.SetPathValue("requestId", "sr1")
	rec = httptest.NewRecorder()
	h.AssignServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestCancelServiceRequestHTTP_Errors(t *testing.T) {
	_, _, _, _, mockSR, _, h := setupManager(t)

	// missing reqID
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/", nil)
	rec := httptest.NewRecorder()
	h.CancelServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// error branch
	mockSR.EXPECT().CancelServiceRequestByID("s1").Return(errors.New("fail"))
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/s1", nil)
	req.SetPathValue("requestId", "s1")
	rec = httptest.NewRecorder()
	h.CancelServiceRequestHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestListAllFeedbackHTTP_Branches(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)

	// error branch
	mockManager.EXPECT().ViewAllFeedback().Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?all=true", nil)
	rec := httptest.NewRecorder()
	h.ListAllFeedbackHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}

	// all=false branch
	req = httptest.NewRequest(http.MethodGet, "/api/v1/feedback?all=false", nil)
	rec = httptest.NewRecorder()
	h.ListAllFeedbackHTTP(rec, managerReq(req))
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rec.Code)
	}
}

func TestGenerateReportHTTP_Error(t *testing.T) {
	_, _, _, _, _, mockManager, h := setupManager(t)
	mockManager.EXPECT().GetHotelReport().Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	rec := httptest.NewRecorder()
	h.GenerateReportHTTP(rec, managerReq(req))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", rec.Code)
	}
}

func TestAddRoomHTTP_Errors(t *testing.T) {
	_, mockRoom, _, _, _, _, h := setupManager(t)

	// invalid body
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", strings.NewReader("bad"))
	rec := httptest.NewRecorder()
	h.AddRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}

	// AddRoom error
	mockRoom.EXPECT().AddRoom(1, "x", 2.0, true, "").Return(errors.New("fail"))
	body := `{"number":1,"type":"x","price":2}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/rooms", strings.NewReader(body))
	rec = httptest.NewRecorder()
	h.AddRoomHTTP(rec, managerReq(req))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rec.Code)
	}
}

func TestRequireManager(t *testing.T) {
	// no role in context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := RequireManager(req)
	if err == nil {
		t.Error("expected error")
	}

	// wrong role
	req = req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Guest"))
	err = RequireManager(req)
	if err == nil {
		t.Error("expected error")
	}

	// manager role
	req = req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager"))
	err = RequireManager(req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

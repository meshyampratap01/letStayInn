package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	logger "github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

func mgrCtx(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager"))
}

func nonMgrCtx(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Guest"))
}

func setupAll(t *testing.T) (
	*gomock.Controller,
	*mocks.MockIRoomService,
	*mocks.MockIBookingService,
	*mocks.MockIUserService,
	*mocks.MockIServiceRequestService,
	*mocks.MockIManagerService,
	*ManagerHandler,
) {
	ctrl := gomock.NewController(t)
	mockRoom := mocks.NewMockIRoomService(ctrl)
	mockBooking := mocks.NewMockIBookingService(ctrl)
	mockUser := mocks.NewMockIUserService(ctrl)
	mockServiceReq := mocks.NewMockIServiceRequestService(ctrl)
	mockManager := mocks.NewMockIManagerService(ctrl)

	logger.Log = zap.NewNop()

	h := NewManagerHandler(mockRoom, mockBooking, mockUser, mockServiceReq, mockManager)
	return ctrl, mockRoom, mockBooking, mockUser, mockServiceReq, mockManager, h
}

// ---------- RequireManager ----------
func TestRequireManager_Behaviour(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// missing role -> unauthorized error
	if err := RequireManager(req); err == nil || !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("expected unauthorized error when role missing")
	}

	// wrong type -> unauthorized
	req = req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, 123))
	if err := RequireManager(req); err == nil || !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("expected unauthorized error for wrong type")
	}

	// wrong role string -> forbidden
	req = req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Guest"))
	if err := RequireManager(req); err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("expected forbidden error for non-manager")
	}

	// correct
	req = req.WithContext(context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager"))
	if err := RequireManager(req); err != nil {
		t.Fatalf("expected no error for manager role, got %v", err)
	}
}

// ---------- UpdateRoomHTTP ----------
func TestUpdateRoomHTTP_AllBranches(t *testing.T) {
	ctrl, mockRoom, _, _, _, _, h := setupAll(t)
	defer ctrl.Finish()

	// Forbidden (not manager)
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/1", bytes.NewBufferString(`{}`))
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 for forbidden, got %d", rec.Code)
		}
	}

	// Invalid JSON
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/1", bytes.NewBufferString("{bad json"))
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "1")
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid json, got %d", rec.Code)
		}
	}

	// Invalid room number via Sscanf -> non-int
	{
		body := `{"choice":1,"type":"X","price":10,"is_available":true,"description":"d"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/abc", bytes.NewBufferString(body))
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "abc")
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid room number, got %d", rec.Code)
		}
	}

	// room number <=0
	{
		body := `{"choice":1,"type":"X","price":10,"is_available":true,"description":"d"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/0", bytes.NewBufferString(body))
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "0")
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for room number zero, got %d", rec.Code)
		}
	}

	// service error
	{
		body := `{"choice":2,"type":"Deluxe","price":150,"is_available":true,"description":"ok"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/101", bytes.NewBufferString(body))
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "101")
		mockRoom.EXPECT().UpdateRoom(101, 2, "Deluxe", 150.0, true, "ok").Return(errors.New("uerr"))
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for update error, got %d", rec.Code)
		}
	}

	// success
	{
		body := `{"choice":2,"type":"Deluxe","price":150,"is_available":false,"description":"ok"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/rooms/102", bytes.NewBufferString(body))
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "102")
		mockRoom.EXPECT().UpdateRoom(102, 2, "Deluxe", 150.0, false, "ok").Return(nil)
		rec := httptest.NewRecorder()
		h.UpdateRoomHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for success update, got %d", rec.Code)
		}
	}
}

// ---------- DeleteRoomHTTP ----------
func TestDeleteRoomHTTP_AllBranches(t *testing.T) {
	ctrl, mockRoom, mockBooking, _, _, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/1", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("roomNum", "1")
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 for forbidden delete, got %d", rec.Code)
		}
	}

	// invalid room number (atoi fail)
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/abc", nil)
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "abc")
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid room number, got %d", rec.Code)
		}
	}

	// booking check error
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/201", nil)
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "201")
		mockBooking.EXPECT().IsRoomBooked(201).Return(false, errors.New("checkerr"))
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 for booking check error, got %d", rec.Code)
		}
	}

	// booked -> cannot delete
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/202", nil)
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "202")
		mockBooking.EXPECT().IsRoomBooked(202).Return(true, nil)
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for booked room, got %d", rec.Code)
		}
	}

	// delete fail
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/203", nil)
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "203")
		mockBooking.EXPECT().IsRoomBooked(203).Return(false, nil)
		mockRoom.EXPECT().DeleteRoom(203).Return(errors.New("delerr"))
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for delete fail, got %d", rec.Code)
		}
	}

	// success
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/204", nil)
		req = mgrCtx(req)
		req.SetPathValue("roomNum", "204")
		mockBooking.EXPECT().IsRoomBooked(204).Return(false, nil)
		mockRoom.EXPECT().DeleteRoom(204).Return(nil)
		rec := httptest.NewRecorder()
		h.DeleteRoomHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for delete success, got %d", rec.Code)
		}
	}
}

// ---------- ListEmployeesHTTP ----------
func TestListEmployeesHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListEmployeesHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 for forbidden list employees, got %d", rec.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
		req = mgrCtx(req)
		mockManager.EXPECT().GetAllEmployees().Return(nil, errors.New("fetcherr"))
		rec := httptest.NewRecorder()
		h.ListEmployeesHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 for service error, got %d", rec.Code)
		}
	}

	// success
	{
		emps := []models.User{
			{ID: "e1", Name: "A", Email: "a@x", Role: models.RoleKitchenStaff, Available: true},
			{ID: "e2", Name: "B", Email: "b@x", Role: models.RoleCleaningStaff, Available: false},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
		req = mgrCtx(req)
		mockManager.EXPECT().GetAllEmployees().Return(emps, nil)
		rec := httptest.NewRecorder()
		h.ListEmployeesHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for list employees success, got %d", rec.Code)
		}
		// quick decode to ensure body JSON path executed
		var wrapper struct {
			Code int                  `json:"code"`
			Data []dto.UserProfileDTO `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&wrapper); err != nil {
			t.Fatalf("failed decode response: %v", err)
		}
		if wrapper.Code != 200 || len(wrapper.Data) != 2 {
			t.Fatalf("unexpected response payload")
		}
	}
}

// ---------- CreateEmployeeHTTP ----------
func TestCreateEmployeeHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, mockUser, _, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(`{}`))
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.CreateEmployeeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 for forbidden create, got %d", rec.Code)
		}
	}

	// invalid json
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString("{bad"))
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.CreateEmployeeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid json, got %d", rec.Code)
		}
	}

	// invalid role string
	{
		body := `{"name":"X","email":"x@x","password":"p","role":"Invalid","available":true}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(body))
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.CreateEmployeeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid role, got %d", rec.Code)
		}
	}

	// create user service error
	{
		body := `{"name":"John","email":"john@test","password":"p","role":"KitchenStaff","available":true}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(body))
		req = mgrCtx(req)
		mockUser.EXPECT().CreateEmployee("John", "john@test", "p", models.RoleKitchenStaff, true).Return(models.User{}, errors.New("createfail"))
		rec := httptest.NewRecorder()
		h.CreateEmployeeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for create fail, got %d", rec.Code)
		}
	}

	// success
	{
		body := `{"name":"John","email":"john@test","password":"p","role":"CleaningStaff","available":false}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBufferString(body))
		req = mgrCtx(req)
		mockUser.EXPECT().CreateEmployee("John", "john@test", "p", models.RoleCleaningStaff, false).Return(models.User{ID: "u123"}, nil)
		rec := httptest.NewRecorder()
		h.CreateEmployeeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 for create success, got %d", rec.Code)
		}
	}
}

// ---------- DeleteEmployeeHTTP ----------
func TestDeleteEmployeeHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/id", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("employeeId", "id")
		rec := httptest.NewRecorder()
		h.DeleteEmployeeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/", nil)
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.DeleteEmployeeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 missing id, got %d", rec.Code)
		}
	}

	// delete error
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/id", nil)
		req = mgrCtx(req)
		req.SetPathValue("employeeId", "id")
		mockManager.EXPECT().DeleteEmployeeByID("id").Return(errors.New("delerr"))
		rec := httptest.NewRecorder()
		h.DeleteEmployeeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 delete error, got %d", rec.Code)
		}
	}

	// success
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/id2", nil)
		req = mgrCtx(req)
		req.SetPathValue("employeeId", "id2")
		mockManager.EXPECT().DeleteEmployeeByID("id2").Return(nil)
		rec := httptest.NewRecorder()
		h.DeleteEmployeeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 delete success, got %d", rec.Code)
		}
	}
}

// ---------- UpdateEmployeeAvailabilityHTTP ----------
func TestUpdateEmployeeAvailabilityHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/email/availability", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("employeeEmail", "x@x")
		rec := httptest.NewRecorder()
		h.UpdateEmployeeAvailabilityHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// missing email
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/employees//availability", nil)
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.UpdateEmployeeAvailabilityHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 missing email, got %d", rec.Code)
		}
	}

	// invalid JSON
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/x@x/availability", bytes.NewBufferString("{bad"))
		req = mgrCtx(req)
		req.SetPathValue("employeeEmail", "x@x")
		rec := httptest.NewRecorder()
		h.UpdateEmployeeAvailabilityHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 invalid json, got %d", rec.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/x@x/availability", bytes.NewBufferString(`{"available":true}`))
		req = mgrCtx(req)
		req.SetPathValue("employeeEmail", "x@x")
		mockManager.EXPECT().UpdateEmployeeAvailability("x@x", true).Return(errors.New("uerr"))
		rec := httptest.NewRecorder()
		h.UpdateEmployeeAvailabilityHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 service error, got %d", rec.Code)
		}
	}

	// success
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/employees/x@x/availability", bytes.NewBufferString(`{"available":false}`))
		req = mgrCtx(req)
		req.SetPathValue("employeeEmail", "x@x")
		mockManager.EXPECT().UpdateEmployeeAvailability("x@x", false).Return(nil)
		rec := httptest.NewRecorder()
		h.UpdateEmployeeAvailabilityHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 success, got %d", rec.Code)
		}
	}
}

// ---------- ListAllBookingsHTTP ----------
func TestListAllBookingsHTTP_AllBranches(t *testing.T) {
	ctrl, _, mockBooking, _, _, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListAllBookingsHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// all != true -> empty bookings (success)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListAllBookingsHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for empty bookings, got %d", rec.Code)
		}
	}

	// all=true but service error
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings?all=true", nil)
		req = mgrCtx(req)
		req = req.WithContext(req.Context()) // keep same
		// ensure query param present
		req.URL.RawQuery = "all=true"
		mockBooking.EXPECT().GetActiveBookings().Return(nil, errors.New("fetcherr"))
		rec := httptest.NewRecorder()
		h.ListAllBookingsHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 for booking fetch error, got %d", rec.Code)
		}
	}

	// success with active bookings
	{
		bs := []models.Booking{
			{ID: "b1", RoomNum: 10},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings?all=true", nil)
		req = mgrCtx(req)
		req.URL.RawQuery = "all=true"
		mockBooking.EXPECT().GetActiveBookings().Return(bs, nil)
		rec := httptest.NewRecorder()
		h.ListAllBookingsHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for bookings success, got %d", rec.Code)
		}
	}
}

// ---------- ListUnassignedServiceRequestsHTTP ----------
func TestListUnassignedServiceRequestsHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, mockServiceReq, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/service-requests", nil)
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListUnassignedServiceRequestsHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/service-requests", nil)
		req = mgrCtx(req)
		mockServiceReq.EXPECT().GetUnassignedServiceRequest().Return(nil, errors.New("serr"))
		rec := httptest.NewRecorder()
		h.ListUnassignedServiceRequestsHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 on service error, got %d", rec.Code)
		}
	}

	// success
	{
		now := time.Now()
		reqdata := []models.ServiceRequest{
			{ID: "sr1", RoomNum: 11, Type: models.ServiceTypeCleaning, Details: "d", Status: models.ServiceStatusPending, AssignedTo: "", IsAssigned: false, CreatedAt: now},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/service-requests", nil)
		req = mgrCtx(req)
		mockServiceReq.EXPECT().GetUnassignedServiceRequest().Return(reqdata, nil)
		rec := httptest.NewRecorder()
		h.ListUnassignedServiceRequestsHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for list unassigned success, got %d", rec.Code)
		}
	}
}

// ---------- AssignServiceRequestHTTP ----------
func TestAssignServiceRequestHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/id/assign", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("requestId", "id")
		rec := httptest.NewRecorder()
		h.AssignServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// missing requestId
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests//assign", bytes.NewBufferString(`{"employee_id":"e1"}`))
		req = mgrCtx(req)
		// ensure empty path value
		req.SetPathValue("requestId", "")
		rec := httptest.NewRecorder()
		h.AssignServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for missing request id, got %d", rec.Code)
		}
	}

	// invalid body
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/r1/assign", bytes.NewBufferString("{bad"))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r1")
		rec := httptest.NewRecorder()
		h.AssignServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid body, got %d", rec.Code)
		}
	}

	// service assign error
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/r2/assign", bytes.NewBufferString(`{"employee_id":"e2"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r2")
		mockManager.EXPECT().AssignServiceRequest("r2", "e2").Return(errors.New("assignfail"))
		rec := httptest.NewRecorder()
		h.AssignServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for assign fail, got %d", rec.Code)
		}
	}

	// success
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests/r3/assign", bytes.NewBufferString(`{"employee_id":"e3"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r3")
		mockManager.EXPECT().AssignServiceRequest("r3", "e3").Return(nil)
		rec := httptest.NewRecorder()
		h.AssignServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 on assign success, got %d", rec.Code)
		}
	}
}

// ---------- UpdateServiceRequestStatusHTTP ----------
func TestUpdateServiceRequestStatusHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, mockServiceReq, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/r/status", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("requestId", "r")
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests//status", bytes.NewBufferString(`{"status":"Pending"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "")
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 missing id, got %d", rec.Code)
		}
	}

	// invalid body
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/r1/status", bytes.NewBufferString("{bad"))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r1")
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 invalid body, got %d", rec.Code)
		}
	}

	// invalid status
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/r2/status", bytes.NewBufferString(`{"status":"NotAStatus"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r2")
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 invalid status, got %d", rec.Code)
		}
	}

	// service update error
	{
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/r3/status", bytes.NewBufferString(`{"status":"Pending"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r3")
		mockServiceReq.EXPECT().UpdateServiceRequestStatus("r3", models.ServiceStatusPending).Return(errors.New("upfail"))
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 update error, got %d", rec.Code)
		}
	}

	// success (test all valid statuses quickly)
	validStatuses := []string{
		string(models.ServiceStatusPending),
		string(models.ServiceStatusInProgress),
		string(models.ServiceStatusDone),
		string(models.ServiceStatusCancelled),
	}
	for i, s := range validStatuses {
		reqID := "ok" + strconv.Itoa(i)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/service-requests/"+reqID+"/status", bytes.NewBufferString(`{"status":"`+s+`"}`))
		req = mgrCtx(req)
		req.SetPathValue("requestId", reqID)
		mockServiceReq.EXPECT().UpdateServiceRequestStatus(reqID, models.ServiceStatus(s)).Return(nil)
		rec := httptest.NewRecorder()
		h.UpdateServiceRequestStatusHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for status %s, got %d", s, rec.Code)
		}
	}
}

// ---------- CancelServiceRequestHTTP ----------
func TestCancelServiceRequestHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, mockServiceReq, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/r", nil)
		req = nonMgrCtx(req)
		req.SetPathValue("requestId", "r")
		rec := httptest.NewRecorder()
		h.CancelServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// missing id
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/", nil)
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.CancelServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 missing id, got %d", rec.Code)
		}
	}

	// cancel error
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/r4", nil)
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r4")
		mockServiceReq.EXPECT().CancelServiceRequestByID("r4").Return(errors.New("cancelerr"))
		rec := httptest.NewRecorder()
		h.CancelServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 cancel error, got %d", rec.Code)
		}
	}

	// success
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/service-requests/r5", nil)
		req = mgrCtx(req)
		req.SetPathValue("requestId", "r5")
		mockServiceReq.EXPECT().CancelServiceRequestByID("r5").Return(nil)
		rec := httptest.NewRecorder()
		h.CancelServiceRequestHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 cancel success, got %d", rec.Code)
		}
	}
}

// ---------- ListAllFeedbackHTTP ----------
func TestListAllFeedbackHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback", nil)
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListAllFeedbackHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// all != true -> empty success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback", nil)
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.ListAllFeedbackHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 empty feedbacks, got %d", rec.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?all=true", nil)
		req = mgrCtx(req)
		req.URL.RawQuery = "all=true"
		mockManager.EXPECT().ViewAllFeedback().Return(nil, errors.New("ferr"))
		rec := httptest.NewRecorder()
		h.ListAllFeedbackHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 feedback fetch error, got %d", rec.Code)
		}
	}

	// success
	{
		f := []models.Feedback{
			{ID: "f1", UserID: "u1", UserName: "U", Message: "ok", RoomNum: 101, BookingID: "b1", Rating: 5, CreatedAt: time.Now()},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/feedback?all=true", nil)
		req = mgrCtx(req)
		req.URL.RawQuery = "all=true"
		mockManager.EXPECT().ViewAllFeedback().Return(f, nil)
		rec := httptest.NewRecorder()
		h.ListAllFeedbackHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 feedback success, got %d", rec.Code)
		}
	}
}

// ---------- GenerateReportHTTP ----------
func TestGenerateReportHTTP_AllBranches(t *testing.T) {
	ctrl, _, _, _, _, mockManager, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.GenerateReportHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// service error
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
		req = mgrCtx(req)
		mockManager.EXPECT().GetHotelReport().Return(nil, errors.New("repoerr"))
		rec := httptest.NewRecorder()
		h.GenerateReportHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 report error, got %d", rec.Code)
		}
	}
}

// ---------- AddRoomHTTP ----------
func TestAddRoomHTTP_AllBranches(t *testing.T) {
	ctrl, mockRoom, _, _, _, _, h := setupAll(t)
	defer ctrl.Finish()

	// forbidden
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBufferString(`{}`))
		req = nonMgrCtx(req)
		rec := httptest.NewRecorder()
		h.AddRoomHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 forbidden, got %d", rec.Code)
		}
	}

	// invalid JSON
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBufferString("{bad"))
		req = mgrCtx(req)
		rec := httptest.NewRecorder()
		h.AddRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 invalid json, got %d", rec.Code)
		}
	}

	// service error
	{
		body := `{"number":301,"type":"Deluxe","price":120,"description":"d"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBufferString(body))
		req = mgrCtx(req)
		mockRoom.EXPECT().AddRoom(301, "Deluxe", 120.0, true, "d").Return(errors.New("adderr"))
		rec := httptest.NewRecorder()
		h.AddRoomHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 add room error, got %d", rec.Code)
		}
	}

	// success
	{
		body := `{"number":302,"type":"Deluxe","price":120,"description":"d2"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBufferString(body))
		req = mgrCtx(req)
		mockRoom.EXPECT().AddRoom(302, "Deluxe", 120.0, true, "d2").Return(nil)
		rec := httptest.NewRecorder()
		h.AddRoomHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 add room success, got %d", rec.Code)
		}
	}
}

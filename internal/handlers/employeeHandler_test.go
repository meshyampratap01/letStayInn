package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupEmployee(t *testing.T) (*gomock.Controller, *mocks.MockIEmployeeService, *EmployeeHandler) {
	ctrl := gomock.NewController(t)
	mockSvc := mocks.NewMockIEmployeeService(ctrl)
	h := NewEmployeeHandler(mockSvc)

	logger.Log = zap.NewNop()

	return ctrl, mockSvc, h
}

func decodeResponse(t *testing.T, body *bytes.Buffer) map[string]any {
	var m map[string]any
	if err := json.NewDecoder(body).Decode(&m); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	return m
}

func TestViewAssignedServiceRequestsHTTP_Success(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return([]models.ServiceRequest{
		{ID: "sr1", RoomNum: 101, Type: models.ServiceTypeCleaning, Status: models.ServiceStatusPending, AssignedTo: "emp1", IsAssigned: true, CreatedAt: time.Now()},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employee/service-requests", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.ViewAssignedServiceRequestsHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
	resp := decodeResponse(t, rec.Body)
	if resp["code"].(float64) != 200 {
		t.Errorf("expected code 200, got %v", resp["code"])
	}
}

func TestViewAssignedServiceRequestsHTTP_InvalidRole(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ViewAssignedServiceRequestsHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestViewAssignedServiceRequestsHTTP_InvalidUserID(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	rec := httptest.NewRecorder()
	h.ViewAssignedServiceRequestsHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestViewAssignedServiceRequestsHTTP_ForbiddenRole(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()
	h.ViewAssignedServiceRequestsHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestViewAssignedServiceRequestsHTTP_ServiceError(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "KitchenStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()
	h.ViewAssignedServiceRequestsHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_Success(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return([]models.ServiceRequest{
		{ID: "sr1", AssignedTo: "emp1"},
	}, nil)
	mockSvc.EXPECT().UpdateServiceRequestStatus("sr1", models.ServiceStatusDone).Return(nil)

	body := `{"status":"Done"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/employee/service-requests/sr1/status", bytes.NewBufferString(body))
	req.SetPathValue("serviceRequestId", "sr1")

	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_MissingID(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{"status":"Done"}`))
	rec := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_InvalidBody(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("{invalid}"))
	req.SetPathValue("serviceRequestId", "sr1")
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_InvalidStatus(t *testing.T) {
	_, _, h := setupEmployee(t)

	body := `{"status":"NotAStatus"}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	req.SetPathValue("serviceRequestId", "sr1")
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "KitchenStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_ServiceFetchError(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return(nil, errors.New("fetch failed"))

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	req.SetPathValue("serviceRequestId", "sr1")
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_NotAssigned(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return([]models.ServiceRequest{}, nil)

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	req.SetPathValue("serviceRequestId", "sr1")
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestUpdateServiceRequestStatusHTTP_UpdateError(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetAssignedServiceRequests("emp1").Return([]models.ServiceRequest{
		{ID: "sr1", AssignedTo: "emp1"},
	}, nil)
	mockSvc.EXPECT().UpdateServiceRequestStatus("sr1", models.ServiceStatusPending).Return(errors.New("fail"))

	body := `{"status":"Pending"}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	req.SetPathValue("serviceRequestId", "sr1")
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "KitchenStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.UpdateServiceRequestStatusHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestToggleAvailabilityHTTP_Success(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().ToggleAvailability("emp1").Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "CleaningStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.ToggleAvailabilityHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestToggleAvailabilityHTTP_Error(t *testing.T) {
	ctrl, mockSvc, h := setupEmployee(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().ToggleAvailability("emp1").Return(errors.New("fail"))

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "KitchenStaff")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.ToggleAvailabilityHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestToggleAvailabilityHTTP_InvalidRole(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	rec := httptest.NewRecorder()

	h.ToggleAvailabilityHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestToggleAvailabilityHTTP_ForbiddenRole(t *testing.T) {
	_, _, h := setupEmployee(t)

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	ctx := context.WithValue(req.Context(), contextkeys.UserRoleKey, "Manager")
	ctx = context.WithValue(ctx, contextkeys.UserIDKey, "emp1")
	rec := httptest.NewRecorder()

	h.ToggleAvailabilityHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

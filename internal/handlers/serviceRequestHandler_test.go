package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
)

func setupServiceRequestHandler(t *testing.T) (*ServiceRequestHandler,
	*mocks.MockIServiceRequestService,
	*mocks.MockIBookingService,
	*gomock.Controller) {

	ctrl := gomock.NewController(t)
	mockSRS := mocks.NewMockIServiceRequestService(ctrl)
	mockBS := mocks.NewMockIBookingService(ctrl)
	h := NewServiceRequestHandler(mockSRS, mockBS)
	return h, mockSRS, mockBS, ctrl
}

func TestSubmitServiceRequestHTTP(t *testing.T) {
	t.Run("missing role in context", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", nil)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
	})

	t.Run("missing userID in context", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
	})

	t.Run("role not Guest", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Manager")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", rr.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewBufferString("{bad json")).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("invalid room number", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		body := `{"room_num":0,"type":"Cleaning","details":"valid details"}`
		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewBufferString(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("details too short", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		body := `{"room_num":101,"type":"Cleaning","details":"bad"}`
		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewBufferString(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("invalid service type", func(t *testing.T) {
		h, _, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		body := `{"room_num":101,"type":"Invalid","details":"valid details"}`
		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewBufferString(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("service request service returns error", func(t *testing.T) {
		h, mockSRS, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		body := `{"room_num":101,"type":"Cleaning","details":"valid details"}`
		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")

		mockSRS.EXPECT().ServiceRequestGetter(gomock.Any(), 101, models.ServiceTypeCleaning, "valid details").
			Return(errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewBufferString(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		h, mockSRS, _, ctrl := setupServiceRequestHandler(t)
		defer ctrl.Finish()

		mockSRS.EXPECT().ServiceRequestGetter(gomock.Any(), 101, models.ServiceTypeCleaning, "Clean room").Return(nil)

		ctx := context.WithValue(context.Background(), contextkeys.UserRoleKey, "Guest")
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, "u1")
		body := []byte(`{"room_num":101,"type":"Cleaning","details":"Clean room"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", bytes.NewReader(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		h.SubmitServiceRequestHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d", rr.Code)
		}

		var wrapper struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    any    `json:"data"`
		}
		if err := json.NewDecoder(rr.Body).Decode(&wrapper); err != nil {
			t.Fatal("failed to decode response")
		}
		if wrapper.Message != "Service request submitted" {
			t.Errorf("expected success message, got %s", wrapper.Message)
		}
	})
}

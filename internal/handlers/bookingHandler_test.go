package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contextKeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setup(t *testing.T) (*gomock.Controller, *mocks.MockIBookingService, *mocks.MockIRoomService, *BookingHandler) {
	ctrl := gomock.NewController(t)
	mockBooking := mocks.NewMockIBookingService(ctrl)
	mockRoom := mocks.NewMockIRoomService(ctrl)
	h := NewBookingHandler(mockBooking, mockRoom)

	logger.Log = zap.NewNop()

	return ctrl, mockBooking, mockRoom, h
}

func TestGetRoomsByRoleHTTP_Manager_AllRooms(t *testing.T) {
	ctrl, _, mockRoom, h := setup(t)
	defer ctrl.Finish()

	mockRoom.EXPECT().GetAllRooms().Return([]models.Room{{Number: 101}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
	rec := httptest.NewRecorder()

	h.GetRoomsByRoleHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp []dto.RoomDTO
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp) != 1 || resp[0].Number != 101 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetRoomsByRoleHTTP_OtherRole_OnlyAvailable(t *testing.T) {
	ctrl, _, mockRoom, h := setup(t)
	defer ctrl.Finish()

	mockRoom.EXPECT().GetAvailableRooms().Return([]models.Room{{Number: 102}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Guest")
	rec := httptest.NewRecorder()

	h.GetRoomsByRoleHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestBookRoomHTTP_Success(t *testing.T) {
	ctrl, mockBooking, _, h := setup(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().BookRoom(gomock.Any(), 101, "2025-08-01", "2025-08-02").Return(nil)

	body := `{"room_number":101,"check_in_date":"2025-08-01","check_out_date":"2025-08-02"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
	rec := httptest.NewRecorder()

	h.BookRoomHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestBookRoomHTTP_InvalidBody(t *testing.T) {
	_, _, _, h := setup(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", strings.NewReader("{invalid}"))
	rec := httptest.NewRecorder()

	h.BookRoomHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestGetBookingsByRoleHTTP_Manager(t *testing.T) {
	ctrl, mockBooking, _, h := setup(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().GetActiveBookings().Return([]models.Booking{{ID: "b1", RoomNum: 201}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
	rec := httptest.NewRecorder()

	h.GetBookingsByRoleHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCancelBookingHTTP_Success(t *testing.T) {
	ctrl, mockBooking, _, h := setup(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().CancelBooking(gomock.Any(), "b1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/b1", nil)
	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
	rec := httptest.NewRecorder()

	h.CancelBookingHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCancelBookingHTTP_InvalidID(t *testing.T) {
	_, _, _, h := setup(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/", nil)
	rec := httptest.NewRecorder()

	h.CancelBookingHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCancelBookingHTTP_ErrorFromService(t *testing.T) {
	ctrl, mockBooking, _, h := setup(t)
	defer ctrl.Finish()

	mockBooking.EXPECT().CancelBooking(gomock.Any(), "b1").Return(errors.New("cannot cancel"))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/b1", nil)
	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
	rec := httptest.NewRecorder()

	h.CancelBookingHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

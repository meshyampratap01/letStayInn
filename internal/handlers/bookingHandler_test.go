package handlers

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	contextKeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
// 	"github.com/meshyampratap01/letStayInn/internal/dto"
// 	"github.com/meshyampratap01/letStayInn/internal/logger"
// 	"github.com/meshyampratap01/letStayInn/internal/mocks"
// 	"github.com/meshyampratap01/letStayInn/internal/models"
// 	gomock "go.uber.org/mock/gomock"
// 	"go.uber.org/zap"
// )

// func setup(t *testing.T) (*gomock.Controller, *mocks.MockIBookingService, *mocks.MockIRoomService, *BookingHandler) {
// 	ctrl := gomock.NewController(t)
// 	mockBooking := mocks.NewMockIBookingService(ctrl)
// 	mockRoom := mocks.NewMockIRoomService(ctrl)
// 	h := NewBookingHandler(mockBooking, mockRoom)
// 	logger.Log = zap.NewNop()
// 	return ctrl, mockBooking, mockRoom, h
// }

// func TestGetRoomsByRoleHTTP_Manager_AllRooms(t *testing.T) {
// 	ctrl, _, mockRoom, h := setup(t)
// 	defer ctrl.Finish()

// 	mockRoom.EXPECT().GetAllRooms().Return([]models.Room{{Number: 101}}, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
// 	rec := httptest.NewRecorder()

// 	h.GetRoomsByRoleHTTP(rec, req.WithContext(ctx))

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf("expected 200, got %d", rec.Code)
// 	}
// 	var wrapper struct {
// 		Data []dto.RoomDTO `json:"data"`
// 	}
// 	_ = json.NewDecoder(rec.Body).Decode(&wrapper)
// 	if len(wrapper.Data) != 1 || wrapper.Data[0].Number != 101 {
// 		t.Errorf("unexpected response: %+v", wrapper.Data)
// 	}
// }

// func TestGetRoomsByRoleHTTP_Manager_AvailableOnly(t *testing.T) {
// 	ctrl, _, mockRoom, h := setup(t)
// 	defer ctrl.Finish()

// 	mockRoom.EXPECT().GetAvailableRooms().Return([]models.Room{{Number: 103}}, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms?available=true", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
// 	rec := httptest.NewRecorder()

// 	h.GetRoomsByRoleHTTP(rec, req.WithContext(ctx))

// 	if rec.Code != http.StatusOK {
// 		t.Fatalf("expected 200, got %d", rec.Code)
// 	}
// }

// func TestGetRoomsByRoleHTTP_Error(t *testing.T) {
// 	ctrl, _, mockRoom, h := setup(t)
// 	defer ctrl.Finish()

// 	mockRoom.EXPECT().GetAvailableRooms().Return(nil, errors.New("db error"))

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Guest")
// 	rec := httptest.NewRecorder()

// 	h.GetRoomsByRoleHTTP(rec, req.WithContext(ctx))

// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf("expected 500, got %d", rec.Code)
// 	}
// }

// func TestBookRoomHTTP_Success(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().BookRoom(gomock.Any(), 101, "2025-08-01", "2025-08-02").Return(nil)

// 	body := `{"room_number":101,"check_in_date":"2025-08-01","check_out_date":"2025-08-02"}`
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
// 	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
// 	rec := httptest.NewRecorder()

// 	h.BookRoomHTTP(rec, req.WithContext(ctx))

// 	if rec.Code != http.StatusCreated {
// 		t.Fatalf("expected 201, got %d", rec.Code)
// 	}
// }

// func TestBookRoomHTTP_InvalidBody(t *testing.T) {
// 	_, _, _, h := setup(t)

// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", strings.NewReader("{invalid}"))
// 	rec := httptest.NewRecorder()

// 	h.BookRoomHTTP(rec, req)
// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf("expected 400, got %d", rec.Code)
// 	}
// }

// func TestBookRoomHTTP_MissingUserID(t *testing.T) {
// 	_, _, _, h := setup(t)

// 	body := `{"room_number":101,"check_in_date":"2025-08-01","check_out_date":"2025-08-02"}`
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
// 	rec := httptest.NewRecorder()

// 	h.BookRoomHTTP(rec, req)

// 	if rec.Code != http.StatusUnauthorized {
// 		t.Fatalf("expected 401, got %d", rec.Code)
// 	}
// }

// func TestBookRoomHTTP_InvalidRequestFields(t *testing.T) {
// 	_, _, _, h := setup(t)

// 	body := `{"room_number":0,"check_in_date":"","check_out_date":""}`
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
// 	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
// 	rec := httptest.NewRecorder()

// 	h.BookRoomHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf("expected 400, got %d", rec.Code)
// 	}
// }

// func TestBookRoomHTTP_ServiceError(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().BookRoom(gomock.Any(), 101, "2025-08-01", "2025-08-02").Return(errors.New("cannot book"))

// 	body := `{"room_number":101,"check_in_date":"2025-08-01","check_out_date":"2025-08-02"}`
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBufferString(body))
// 	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
// 	rec := httptest.NewRecorder()

// 	h.BookRoomHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf("expected 400, got %d", rec.Code)
// 	}
// }

// func TestGetBookingsByRoleHTTP_Manager(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().GetActiveBookings().Return([]models.Booking{{ID: "b1", RoomNum: 201}}, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
// 	rec := httptest.NewRecorder()

// 	h.GetBookingsByRoleHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusOK {
// 		t.Fatalf("expected 200, got %d", rec.Code)
// 	}
// }

// func TestGetBookingsByRoleHTTP_NonManager(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().GetUserActiveBookings(gomock.Any()).Return([]models.Booking{{ID: "b2"}}, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Guest")
// 	rec := httptest.NewRecorder()

// 	h.GetBookingsByRoleHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusOK {
// 		t.Fatalf("expected 200, got %d", rec.Code)
// 	}
// }

// func TestGetBookingsByRoleHTTP_Error(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().GetActiveBookings().Return(nil, errors.New("db error"))

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings", nil)
// 	ctx := context.WithValue(req.Context(), contextKeys.UserRoleKey, "Manager")
// 	rec := httptest.NewRecorder()

// 	h.GetBookingsByRoleHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusInternalServerError {
// 		t.Fatalf("expected 500, got %d", rec.Code)
// 	}
// }

// // func TestCancelBookingHTTP_Success(t *testing.T) {
// // 	ctrl, mockBooking, _, h := setup(t)
// // 	defer ctrl.Finish()

// // 	mockBooking.EXPECT().CancelBooking(gomock.Any(), "b1").Return(nil)

// // 	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/b1", nil)
// // 	req.SetPathValue("id", "b1")
// // 	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
// // 	rec := httptest.NewRecorder()

// // 	h.CancelBookingHTTP(rec, req.WithContext(ctx))
// // 	if rec.Code != http.StatusOK {
// // 		t.Fatalf("expected 200, got %d", rec.Code)
// // 	}
// // }

// func TestCancelBookingHTTP_InvalidID(t *testing.T) {
// 	_, _, _, h := setup(t)

// 	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/", nil)
// 	rec := httptest.NewRecorder()

// 	h.CancelBookingHTTP(rec, req)
// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf("expected 400, got %d", rec.Code)
// 	}
// }

// func TestCancelBookingHTTP_NoUserID(t *testing.T) {
// 	_, _, _, h := setup(t)

// 	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/b1", nil)
// 	req.SetPathValue("id", "b1")
// 	rec := httptest.NewRecorder()

// 	h.CancelBookingHTTP(rec, req)
// 	if rec.Code != http.StatusUnauthorized {
// 		t.Fatalf("expected 401, got %d", rec.Code)
// 	}
// }

// func TestCancelBookingHTTP_ErrorFromService(t *testing.T) {
// 	ctrl, mockBooking, _, h := setup(t)
// 	defer ctrl.Finish()

// 	mockBooking.EXPECT().CancelBooking(gomock.Any(), "b1").Return(errors.New("cannot cancel"))

// 	req := httptest.NewRequest(http.MethodDelete, "/api/v1/bookings/b1", nil)
// 	req.SetPathValue("id", "b1")
// 	ctx := context.WithValue(req.Context(), contextKeys.UserIDKey, "u1")
// 	rec := httptest.NewRecorder()

// 	h.CancelBookingHTTP(rec, req.WithContext(ctx))
// 	if rec.Code != http.StatusBadRequest {
// 		t.Fatalf("expected 400, got %d", rec.Code)
// 	}
// }

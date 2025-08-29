package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meshyampratap01/letStayInn/internal/logger"
	"go.uber.org/zap"
	
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
)

type BookingHandler struct {
	bookingService bookingService.IBookingService
	roomService    roomService.IRoomService
}

func NewBookingHandler(
	bookingService bookingService.IBookingService,
	roomService roomService.IRoomService,
) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		roomService:    roomService,
	}
}
// GET /api/v1/rooms (returns all rooms for manager, only available rooms for others)
func (h *BookingHandler) GetRoomsByRoleHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	roleStr, _ := roleVal.(string)
	var rooms []models.Room
	var err error
	if roleStr == "Manager" {
		available := r.URL.Query().Get("available")
		if strings.ToLower(available) == "true" {
			rooms, err = h.roomService.GetAvailableRooms()
		} else {
			rooms, err = h.roomService.GetAllRooms()
		}
	} else {
		rooms, err = h.roomService.GetAvailableRooms()
	}
	if err != nil {
		logger.Log.Error("Failed to fetch rooms", zap.Error(err))
		http.Error(w, "Failed to fetch rooms", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Rooms fetched", zap.String("role", roleStr), zap.Int("count", len(rooms)))
	resp := make([]dto.RoomDTO, 0, len(rooms))
	for _, r := range rooms {
		resp = append(resp, dto.RoomDTO{
			Number:      r.Number,
			Type:        string(r.Type),
			Price:       r.Price,
			IsAvailable: r.IsAvailable,
			Description: r.Description,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /api/v1/bookings
func (h *BookingHandler) BookRoomHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomNumber   int    `json:"room_number"`
		CheckInDate  string `json:"check_in_date"`
		CheckOutDate string `json:"check_out_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid booking request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userIDVal := ctx.Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	if req.RoomNumber <= 0 || req.CheckInDate == "" || req.CheckOutDate == "" {
		logger.Log.Warn("Invalid booking request", zap.Any("request", req))
		http.Error(w, "Invalid booking request", http.StatusBadRequest)
		return
	}
	err := h.bookingService.BookRoom(ctx, req.RoomNumber, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		logger.Log.Error("Failed to book room", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Room booked successfully", zap.String("userID", userID), zap.Int("roomNumber", req.RoomNumber))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room booked successfully"})
}

// GET /api/v1/bookings (returns all bookings for manager, user bookings for others)
func (h *BookingHandler) GetBookingsByRoleHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	roleStr, _ := roleVal.(string)
	var bookings []models.Booking
	var err error
	if roleStr == "Manager" {
		bookings, err = h.bookingService.GetActiveBookings()
	} else {
		bookings, err = h.bookingService.GetUserActiveBookings(ctx)
	}
	if err != nil {
		logger.Log.Error("Failed to fetch bookings", zap.Error(err))
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Bookings fetched", zap.String("role", roleStr), zap.Int("count", len(bookings)))
	resp := make([]dto.BookingDTO, 0, len(bookings))
	for _, b := range bookings {
		resp = append(resp, dto.BookingDTO{
			ID:         b.ID,
			RoomNumber: b.RoomNum,
			Status:     b.Status,
			CheckIn:    b.CheckIn.Format("2006-01-02"),
			CheckOut:   b.CheckOut.Format("2006-01-02"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/v1/bookings/{id}
func (h *BookingHandler) CancelBookingHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract booking ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/bookings/")
	if idStr == "" {
		http.Error(w, "Booking ID required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userIDVal := ctx.Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	bookingID := strings.TrimSpace(idStr)
	if bookingID == "" {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}
	err := h.bookingService.CancelBooking(ctx, bookingID)
	if err != nil {
		logger.Log.Error("Failed to cancel booking", zap.Error(err), zap.String("userID", userID), zap.String("bookingID", bookingID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Booking cancelled successfully", zap.String("userID", userID), zap.String("bookingID", bookingID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled successfully"})
}


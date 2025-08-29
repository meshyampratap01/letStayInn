package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	serviceRequest "github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

// POST /api/v1/service-requests
func (s *ServiceRequestHandler) SubmitServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid role in context", zap.Any("context", r.Context()))
		http.Error(w, "Unauthorized: invalid role in context", http.StatusUnauthorized)
		return
	}
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid user ID in context", zap.Any("context", r.Context()))
		http.Error(w, "Unauthorized: invalid user ID in context", http.StatusUnauthorized)
		return
	}
	if role != "Guest" {
		logger.Log.Warn("Forbidden: guest access required", zap.String("role", role))
		http.Error(w, "Forbidden: guest access required", http.StatusForbidden)
		return
	}
	var req struct {
		RoomNum int    `json:"room_num"`
		Type    string `json:"type"`
		Details string `json:"details"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for service request creation", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.RoomNum <= 0 {
		logger.Log.Warn("Invalid room number for service request", zap.Int("room_num", req.RoomNum))
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}
	if len(req.Details) < 5 {
		logger.Log.Warn("Service request details too short", zap.String("details", req.Details))
		http.Error(w, "Details must be at least 5 characters", http.StatusBadRequest)
		return
	}
	var serviceType models.ServiceType
	switch req.Type {
	case string(models.ServiceTypeCleaning):
		serviceType = models.ServiceTypeCleaning
	case string(models.ServiceTypeFood):
		serviceType = models.ServiceTypeFood
	default:
		logger.Log.Warn("Invalid service type for service request", zap.String("type", req.Type))
		http.Error(w, "Invalid service type", http.StatusBadRequest)
		return
	}
	ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
	err := s.ServiceRequestService.ServiceRequestGetter(ctx, req.RoomNum, serviceType, req.Details)
	if err != nil {
		logger.Log.Error("Failed to create service request", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Service request submitted", zap.String("userID", userID), zap.Int("roomNum", req.RoomNum), zap.String("serviceType", req.Type))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service request submitted"})
}

type ServiceRequestHandler struct {
	ServiceRequestService serviceRequest.IServiceRequestService
	BookingService        bookingService.IBookingService
}

func NewServiceRequestHandler(srs serviceRequest.IServiceRequestService, bs bookingService.IBookingService) *ServiceRequestHandler {
	return &ServiceRequestHandler{
		ServiceRequestService: srs,
		BookingService:        bs,
	}
}

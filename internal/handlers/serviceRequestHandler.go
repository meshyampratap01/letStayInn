package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/response"

	"go.uber.org/zap"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	serviceRequest "github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

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

// POST /api/v1/service-requests
func (s *ServiceRequestHandler) SubmitServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid role in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "Unauthorized: invalid role in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid user ID in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "Unauthorized: invalid user ID in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if role != "Guest" {
		logger.Log.Warn("Forbidden: guest access required", zap.String("role", role))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, "Forbidden: guest access required")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		RoomNum int    `json:"room_num"`
		Type    string `json:"type"`
		Details string `json:"details"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for service request creation", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if req.RoomNum <= 0 {
		logger.Log.Warn("Invalid room number for service request", zap.Int("room_num", req.RoomNum))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid room number")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if len(req.Details) < 5 {
		logger.Log.Warn("Service request details too short", zap.String("details", req.Details))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Details must be at least 5 characters")
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid service type")
		json.NewEncoder(w).Encode(resp)
		return
	}
	ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
	err := s.ServiceRequestService.ServiceRequestGetter(ctx, req.RoomNum, serviceType, req.Details)
	if err != nil {
		logger.Log.Error("Failed to create service request", zap.Error(err), zap.String("userID", userID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Service request submitted", zap.String("userID", userID), zap.Int("roomNum", req.RoomNum), zap.String("serviceType", req.Type))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := response.NewSuccessResponse(http.StatusCreated, "Service request submitted", nil)
	json.NewEncoder(w).Encode(resp)
}

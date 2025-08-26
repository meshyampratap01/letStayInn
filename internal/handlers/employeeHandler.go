package handlers

import (
	"encoding/json"
	"net/http"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
	"go.uber.org/zap"
)

// GET /api/v1/employee/service-requests
func (eh *EmployeeHandler) ViewAssignedServiceRequestsHTTP(w http.ResponseWriter, r *http.Request) {
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
	if role != "KitchenStaff" && role != "CleaningStaff" {
		logger.Log.Warn("Forbidden: employee access required", zap.String("role", role))
		http.Error(w, "Forbidden: employee access required", http.StatusForbidden)
		return
	}
	requests, err := eh.employeeService.GetAssignedServiceRequests(userID)
	if err != nil {
		logger.Log.Error("Failed to fetch assigned service requests", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Assigned service requests fetched", zap.String("userID", userID), zap.Int("count", len(requests)))
	resp := make([]dto.ServiceRequestDTO, 0, len(requests))
	for _, sr := range requests {
		resp = append(resp, dto.ServiceRequestDTO{
			ID:         sr.ID,
			RoomNum:    sr.RoomNum,
			Type:       string(sr.Type),
			Details:    sr.Details,
			Status:     string(sr.Status),
			EmployeeID: sr.AssignedTo,
			IsAssigned: sr.IsAssigned,
			CreatedAt:  sr.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PUT /api/v1/employee/service-requests/{serviceRequestId}/status
func (eh *EmployeeHandler) UpdateServiceRequestStatusHTTP(w http.ResponseWriter, r *http.Request) {
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
	if role != "KitchenStaff" && role != "CleaningStaff" {
		logger.Log.Warn("Forbidden: employee access required", zap.String("role", role))
		http.Error(w, "Forbidden: employee access required", http.StatusForbidden)
		return
	}
	reqID := r.PathValue("serviceRequestId")
	if reqID == "" {
		logger.Log.Error("Missing service request id for status update")
		http.Error(w, "Missing service request id", http.StatusBadRequest)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for service request status update", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var newStatus models.ServiceStatus
	switch req.Status {
	case string(models.ServiceStatusPending), string(models.ServiceStatusInProgress), string(models.ServiceStatusDone):
		newStatus = models.ServiceStatus(req.Status)
	default:
		logger.Log.Warn("Invalid status for service request", zap.String("status", req.Status))
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}
	assignedRequests, err := eh.employeeService.GetAssignedServiceRequests(userID)
	if err != nil {
		logger.Log.Error("Failed to fetch assigned service requests for status update", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	found := false
	for _, sr := range assignedRequests {
		if sr.ID == reqID {
			found = true
			break
		}
	}
	if !found {
		logger.Log.Warn("Service request not assigned to employee", zap.String("userID", userID), zap.String("serviceRequestId", reqID))
		http.Error(w, "Service request not assigned to you", http.StatusForbidden)
		return
	}
	err = eh.employeeService.UpdateServiceRequestStatus(reqID, newStatus)
	if err != nil {
		logger.Log.Error("Failed to update service request status", zap.Error(err), zap.String("serviceRequestId", reqID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Service request status updated", zap.String("serviceRequestId", reqID), zap.String("userID", userID), zap.String("status", req.Status))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service request status updated"})
}

// PUT /api/v1/employee/availability
func (eh *EmployeeHandler) ToggleAvailabilityHTTP(w http.ResponseWriter, r *http.Request) {
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
	if role != "KitchenStaff" && role != "CleaningStaff" {
		logger.Log.Warn("Forbidden: employee access required", zap.String("role", role))
		http.Error(w, "Forbidden: employee access required", http.StatusForbidden)
		return
	}
	err := eh.employeeService.ToggleAvailability(userID)
	if err != nil {
		logger.Log.Error("Failed to toggle availability", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Employee availability toggled", zap.String("userID", userID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Availability toggled"})
}

type EmployeeHandler struct {
	employeeService employeeService.IEmployeeService
}

func NewEmployeeHandler(employeeService employeeService.IEmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

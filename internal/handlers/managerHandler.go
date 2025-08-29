package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meshyampratap01/letStayInn/internal/logger"
	"go.uber.org/zap"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	managerservice "github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

type ManagerHandler struct {
	roomService           roomService.IRoomService
	bookingService        bookingService.IBookingService
	userService           userService.IUserService
	serviceRequestService servicerequest.IServiceRequestService
	managerService        managerservice.IManagerService
}

// --- MANAGER HTTP HANDLERS ---

// PUT /api/v1/rooms/{roomNum}
func (h *ManagerHandler) UpdateRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateRoom", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	roomNumStr := r.PathValue("roomNum")
	var req struct {
		Choice      int     `json:"choice"`
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		IsAvailable bool    `json:"is_available"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for UpdateRoom", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	number := 0
	_, err := fmt.Sscanf(roomNumStr, "%d", &number)
	if err != nil || number <= 0 {
		logger.Log.Warn("Invalid room number for UpdateRoom", zap.String("roomId", roomNumStr))
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}
	err = h.roomService.UpdateRoom(number, req.Choice, req.Type, req.Price, req.IsAvailable, req.Description)
	if err != nil {
		logger.Log.Error("Failed to update room", zap.Error(err), zap.Int("roomId", number))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Room updated successfully", zap.Int("roomId", number))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room updated successfully"})
}

// DELETE /api/v1/rooms/{roomNum}
func (h *ManagerHandler) DeleteRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for DeleteRoom", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	roomNumStr := r.PathValue("roomNum")
	number, err := strconv.Atoi(roomNumStr)
	if err != nil || number <= 0 {
		logger.Log.Warn("Invalid room number for DeleteRoom", zap.String("roomId", roomNumStr), zap.Int("roomNum", number))
		http.Error(w, "Invalid room number. Must be a positive integer.", http.StatusBadRequest)
		return
	}
	booked, err := h.bookingService.IsRoomBooked(number)
	if err != nil {
		logger.Log.Error("Error checking room booking status", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if booked {
		http.Error(w, "Cannot delete room: it is currently booked", http.StatusBadRequest)
		return
	}
	err = h.roomService.DeleteRoom(number)
	if err != nil {
		logger.Log.Error("Failed to delete room", zap.Error(err), zap.Int("roomId", number))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Room deleted successfully", zap.Int("roomId", number))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room deleted successfully"})
}

// GET /api/v1/employees
func (h *ManagerHandler) ListEmployeesHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListEmployees", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	employees, err := h.managerService.GetAllEmployees()
	if err != nil {
		logger.Log.Error("Error fetching employees", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dtos := make([]dto.UserProfileDTO, 0, len(employees))
	for _, emp := range employees {
		dtos = append(dtos, dto.UserProfileDTO{
			ID:        emp.ID,
			Name:      emp.Name,
			Email:     emp.Email,
			Role:      emp.Role.String(),
			Available: emp.Available,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos)
}

// POST /api/v1/employees
func (h *ManagerHandler) CreateEmployeeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for CreateEmployee", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"`
		Available bool   `json:"available"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for CreateEmployee", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Convert role string to models.Role
	var empRole models.Role
	switch req.Role {
	case "KitchenStaff":
		empRole = models.RoleKitchenStaff
	case "CleaningStaff":
		empRole = models.RoleCleaningStaff
	default:
		http.Error(w, "Invalid role for employee", http.StatusBadRequest)
		return
	}
	emp, err := h.userService.CreateEmployee(req.Name, req.Email, req.Password, empRole, req.Available)
	if err != nil {
		logger.Log.Error("Failed to create employee", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Employee created successfully", zap.String("employeeId", emp.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// DELETE /api/v1/employees/{employeeId}
func (h *ManagerHandler) DeleteEmployeeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for DeleteEmployee", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	employeeId := r.PathValue("employeeId")
	if employeeId == "" {
		http.Error(w, "Missing employee id", http.StatusBadRequest)
		return
	}
	err := h.managerService.DeleteEmployeeByEmail(employeeId)
	if err != nil {
		logger.Log.Error("Failed to delete employee", zap.Error(err), zap.String("employeeId", employeeId))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Employee deleted successfully", zap.String("employeeId", employeeId))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}

// PUT /api/v1/employees/{employeeId}/availability
func (h *ManagerHandler) UpdateEmployeeAvailabilityHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateEmployeeAvailability", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	employeeId := r.PathValue("employeeId")
	if employeeId == "" {
		http.Error(w, "Missing employee id", http.StatusBadRequest)
		return
	}
	var req struct {
		Available bool `json:"available"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for UpdateEmployeeAvailability", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := h.managerService.UpdateEmployeeAvailability(employeeId, req.Available)
	if err != nil {
		logger.Log.Error("Failed to update employee availability", zap.Error(err), zap.String("employeeId", employeeId))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Employee availability updated", zap.String("employeeId", employeeId))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee availability updated"})
}

// GET /api/v1/bookings?all=true
func (h *ManagerHandler) ListAllBookingsHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListAllBookings", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	all := r.URL.Query().Get("all")
	var bookings []models.Booking
	var err error
	if all == "true" {
		bookings, err = h.bookingService.GetActiveBookings()
	} else {
		bookings = []models.Booking{} // or error
	}
	if err != nil {
		logger.Log.Error("Error fetching bookings", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// GET /api/v1/service-requests
func (h *ManagerHandler) ListUnassignedServiceRequestsHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListUnassignedServiceRequests", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	unassigned := r.URL.Query().Get("unassigned")
	var reqs []models.ServiceRequest
	var err error
	if unassigned == "true" {
		reqs, err = h.serviceRequestService.GetUnassignedServiceRequest()
	} else {
		reqs = []models.ServiceRequest{} 
	}
	if err != nil {
		logger.Log.Error("Error fetching service requests", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqs)
}

// POST /api/v1/service-requests/{requestId}/assign
func (h *ManagerHandler) AssignServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for AssignServiceRequest", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		http.Error(w, "Missing service request id", http.StatusBadRequest)
		return
	}
	var req struct {
		EmployeeID string `json:"employee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for AssignServiceRequest", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := h.managerService.AssignServiceRequest(reqID, req.EmployeeID)
	if err != nil {
		logger.Log.Error("Failed to assign service request", zap.Error(err), zap.String("requestId", reqID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Service request assigned", zap.String("requestId", reqID), zap.String("employeeId", req.EmployeeID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service request assigned"})
}

// PUT /api/v1/service-requests/{requestId}/status
func (h *ManagerHandler) UpdateServiceRequestStatusHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateServiceRequestStatus", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		http.Error(w, "Missing service request id", http.StatusBadRequest)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for UpdateServiceRequestStatus", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate status
	var status models.ServiceStatus
	switch req.Status {
	case string(models.ServiceStatusPending), string(models.ServiceStatusInProgress), string(models.ServiceStatusDone), string(models.ServiceStatusCancelled):
		status = models.ServiceStatus(req.Status)
	default:
		http.Error(w, "Invalid service status", http.StatusBadRequest)
		return
	}
	err := h.serviceRequestService.UpdateServiceRequestStatus(reqID, status)
	if err != nil {
		logger.Log.Error("Failed to update service request status", zap.Error(err), zap.String("requestId", reqID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Service request status updated", zap.String("requestId", reqID), zap.String("status", req.Status))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service request status updated"})
}

// DELETE /api/v1/service-requests/{requestId}
func (h *ManagerHandler) CancelServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for CancelServiceRequest", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		http.Error(w, "Missing service request id", http.StatusBadRequest)
		return
	}
	err := h.serviceRequestService.CancelServiceRequestByID(reqID)
	if err != nil {
		logger.Log.Error("Failed to cancel service request", zap.Error(err), zap.String("requestId", reqID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("Service request cancelled", zap.String("requestId", reqID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Service request cancelled"})
}

// GET /api/v1/feedback?all=true
func (h *ManagerHandler) ListAllFeedbackHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListAllFeedback", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	all := r.URL.Query().Get("all")
	var feedbacks []models.Feedback
	var err error
	if all == "true" {
		feedbacks, err = h.managerService.ViewAllFeedback()
	} else {
		feedbacks = []models.Feedback{} // or error
	}
	if err != nil {
		logger.Log.Error("Error fetching feedback", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]dto.FeedbackDTO, 0, len(feedbacks))
	for _, f := range feedbacks {
		resp = append(resp, dto.FeedbackDTO{
			ID:        f.ID,
			UserID:    f.UserID,
			UserName:  f.UserName,
			Message:   f.Message,
			CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
			RoomNum:   f.RoomNum,
			BookingID: f.BookingID,
			Rating:    f.Rating,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/report
func (h *ManagerHandler) GenerateReportHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for GenerateReport", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	report, err := h.managerService.GetHotelReport()
	if err != nil {
		logger.Log.Error("Error generating report", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// POST /api/v1/rooms (manager only)
func (h *ManagerHandler) AddRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req struct {
		Number      int     `json:"number"`
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.roomService.AddRoom(req.Number, req.Type, req.Price, true, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room created successfully"})
}

func NewManagerHandler(rs roomService.IRoomService, bs bookingService.IBookingService, us userService.IUserService, srs servicerequest.IServiceRequestService, ms managerservice.IManagerService) *ManagerHandler {
	return &ManagerHandler{
		roomService:           rs,
		bookingService:        bs,
		userService:           us,
		serviceRequestService: srs,
		managerService:        ms,
	}
}

// RequireManager checks if the current user is a manager, returns error if not
func RequireManager(r *http.Request) error {
	roleVal := r.Context().Value(contextkeys.UserRoleKey)
	roleStr, ok := roleVal.(string)
	if !ok {
		return errors.New("unauthorized: invalid role in context")
	}
	if roleStr != "Manager" {
		return errors.New("forbidden: manager access required")
	}
	return nil
}

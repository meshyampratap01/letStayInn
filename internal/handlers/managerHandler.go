package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meshyampratap01/letStayInn/internal/response"

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

// --- MANAGER HTTP HANDLERS ---

// PUT /api/v1/rooms/{roomNum}
func (h *ManagerHandler) UpdateRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateRoom", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	number := 0
	_, err := fmt.Sscanf(roomNumStr, "%d", &number)
	if err != nil || number <= 0 {
		logger.Log.Warn("Invalid room number for UpdateRoom", zap.String("roomId", roomNumStr))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid room number")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = h.roomService.UpdateRoom(number, req.Choice, req.Type, req.Price, req.IsAvailable, req.Description)
	if err != nil {
		logger.Log.Error("Failed to update room", zap.Error(err), zap.Int("roomId", number))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Room updated successfully", zap.Int("roomId", number))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Room updated successfully", nil)
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/v1/rooms/{roomNum}
func (h *ManagerHandler) DeleteRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for DeleteRoom", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	roomNumStr := r.PathValue("roomNum")
	number, err := strconv.Atoi(roomNumStr)
	if err != nil || number <= 0 {
		logger.Log.Warn("Invalid room number for DeleteRoom", zap.String("roomId", roomNumStr), zap.Int("roomNum", number))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid room number. Must be a positive integer.")
		json.NewEncoder(w).Encode(resp)
		return
	}
	booked, err := h.bookingService.IsRoomBooked(number)
	if err != nil {
		logger.Log.Error("Error checking room booking status", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	if booked {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Cannot delete room: it is currently booked")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = h.roomService.DeleteRoom(number)
	if err != nil {
		logger.Log.Error("Failed to delete room", zap.Error(err), zap.Int("roomId", number))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Room deleted successfully", zap.Int("roomId", number))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Room deleted successfully", nil)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/employees
func (h *ManagerHandler) ListEmployeesHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListEmployees", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	employees, err := h.managerService.GetAllEmployees()
	if err != nil {
		logger.Log.Error("Error fetching employees", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
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
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Employees fetched successfully", dtos)
	json.NewEncoder(w).Encode(resp)
}

// POST /api/v1/employees
func (h *ManagerHandler) CreateEmployeeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for CreateEmployee", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid role for employee")
		json.NewEncoder(w).Encode(resp)
		return
	}
	emp, err := h.userService.CreateEmployee(req.Name, req.Email, req.Password, empRole, req.Available)
	if err != nil {
		logger.Log.Error("Failed to create employee", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Employee created successfully", zap.String("employeeId", emp.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := response.NewSuccessResponse(http.StatusCreated, "Employee created successfully", emp)
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/v1/employees/{employeeId}
func (h *ManagerHandler) DeleteEmployeeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for DeleteEmployee", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	employeeId := r.PathValue("employeeId")
	if employeeId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Missing employee id")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := h.managerService.DeleteEmployeeByID(employeeId) // param is now id
	if err != nil {
		logger.Log.Error("Failed to delete employee", zap.Error(err), zap.String("employeeId", employeeId))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Employee deleted successfully", zap.String("employeeId", employeeId))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Employee deleted successfully", nil)
	json.NewEncoder(w).Encode(resp)
}

// PUT /api/v1/employees/{employeeEmail}/availability
func (h *ManagerHandler) UpdateEmployeeAvailabilityHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateEmployeeAvailability", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	employeeEmail := r.PathValue("employeeEmail")
	if employeeEmail == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Missing employee id")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		Available bool `json:"available"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for UpdateEmployeeAvailability", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := h.managerService.UpdateEmployeeAvailability(employeeEmail, req.Available)
	if err != nil {
		logger.Log.Error("Failed to update employee availability", zap.Error(err), zap.String("employeeEmail", employeeEmail))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Employee availability updated", zap.String("employeeEmail", employeeEmail))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Employee availability updated", nil)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/bookings?all=true
func (h *ManagerHandler) ListAllBookingsHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListAllBookings", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Bookings fetched successfully", bookings)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/service-requests
func (h *ManagerHandler) ListUnassignedServiceRequestsHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListUnassignedServiceRequests", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	reqs, err := h.serviceRequestService.GetUnassignedServiceRequest()
	if err != nil {
		logger.Log.Error("Error fetching service requests", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	dtos := make([]dto.ServiceRequestDTO, 0, len(reqs))
	for _, sr := range reqs {
		dtos = append(dtos, dto.ServiceRequestDTO{
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
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Service requests fetched successfully", dtos)
	json.NewEncoder(w).Encode(resp)
}

// POST /api/v1/service-requests/{requestId}/assign
func (h *ManagerHandler) AssignServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for AssignServiceRequest", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Missing service request id")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		EmployeeID string `json:"employee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for AssignServiceRequest", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := h.managerService.AssignServiceRequest(reqID, req.EmployeeID)
	if err != nil {
		logger.Log.Error("Failed to assign service request", zap.Error(err), zap.String("requestId", reqID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Service request assigned", zap.String("requestId", reqID), zap.String("employeeId", req.EmployeeID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Service request assigned", nil)
	json.NewEncoder(w).Encode(resp)
}

// PUT /api/v1/service-requests/{requestId}/status
func (h *ManagerHandler) UpdateServiceRequestStatusHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for UpdateServiceRequestStatus", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Missing service request id")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for UpdateServiceRequestStatus", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	// Validate status
	var status models.ServiceStatus
	switch req.Status {
	case string(models.ServiceStatusPending), string(models.ServiceStatusInProgress), string(models.ServiceStatusDone), string(models.ServiceStatusCancelled):
		status = models.ServiceStatus(req.Status)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid service status")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := h.serviceRequestService.UpdateServiceRequestStatus(reqID, status)
	if err != nil {
		logger.Log.Error("Failed to update service request status", zap.Error(err), zap.String("requestId", reqID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Service request status updated", zap.String("requestId", reqID), zap.String("status", req.Status))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Service request status updated", nil)
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/v1/service-requests/{requestId}
func (h *ManagerHandler) CancelServiceRequestHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for CancelServiceRequest", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	reqID := r.PathValue("requestId")
	if reqID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Missing service request id")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := h.serviceRequestService.CancelServiceRequestByID(reqID)
	if err != nil {
		logger.Log.Error("Failed to cancel service request", zap.Error(err), zap.String("requestId", reqID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Service request cancelled", zap.String("requestId", reqID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Service request cancelled", nil)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/feedback?all=true
func (h *ManagerHandler) ListAllFeedbackHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for ListAllFeedback", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	dtos := make([]dto.FeedbackDTO, 0, len(feedbacks))
	for _, f := range feedbacks {
		dtos = append(dtos, dto.FeedbackDTO{
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
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Feedback fetched successfully", dtos)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/report
func (h *ManagerHandler) GenerateReportHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		logger.Log.Warn("Manager role required for GenerateReport", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	report, err := h.managerService.GetHotelReport()
	if err != nil {
		logger.Log.Error("Error generating report", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Report generated successfully", report)
	json.NewEncoder(w).Encode(resp)
}

// POST /api/v1/rooms (manager only)
func (h *ManagerHandler) AddRoomHTTP(w http.ResponseWriter, r *http.Request) {
	if err := RequireManager(r); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req struct {
		Number      int     `json:"number"`
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}

	err := h.roomService.AddRoom(req.Number, req.Type, req.Price, true, req.Description)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := response.NewSuccessResponse(http.StatusCreated, "Room created successfully", nil)
	json.NewEncoder(w).Encode(resp)
}

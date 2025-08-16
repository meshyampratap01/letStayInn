package servicerequest

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type ServiceRequestService struct {
	bookingRepo        bookingRepository.BookingRepository
	serviceRequestRepo serviceRequestRepository.ServiceRequestRepository
}

func NewServiceRequestService(bookingRepo bookingRepository.BookingRepository, serviceRequestRepo serviceRequestRepository.ServiceRequestRepository) IServiceRequestService {
	return &ServiceRequestService{
		bookingRepo:        bookingRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (s *ServiceRequestService) ServiceRequestGetter(ctx context.Context, roomNum int, reqType models.ServiceType, details string) error {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return fmt.Errorf("failed to load bookings: %w", err)
	}

	hasValidBooking := false
	var bid string
	for _, b := range bookings {
		if b.UserID == ctx.Value(contextkeys.UserIDKey) && b.RoomNum == roomNum &&
			b.Status != models.BookingStatusCancelled {
			bid = b.ID
			hasValidBooking = true
			break
		}
	}
	if !hasValidBooking {
		return fmt.Errorf("you don't have any active or completed booking for room %d", roomNum)
	}

	requests, err := s.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return fmt.Errorf("failed to load service requests: %w", err)
	}

	var latestRequest *models.ServiceRequest
	for _, r := range requests {
		if r.UserID == ctx.Value(contextkeys.UserIDKey) && r.RoomNum == roomNum && r.Type == reqType {
			if latestRequest == nil || r.CreatedAt.After(latestRequest.CreatedAt) {
				latestRequest = &r
			}
		}
	}

	if latestRequest != nil {
		if latestRequest.Status != models.ServiceStatusDone && latestRequest.Status != models.ServiceStatusCancelled {
			return fmt.Errorf("your previous %s request for this room is still in progress", reqType)
		}
	}

	newRequest := models.ServiceRequest{
		ID:        utils.NewUUID(),
		UserID:    ctx.Value(contextkeys.UserIDKey).(string),
		RoomNum:   roomNum,
		BookingID: bid,
		Type:      reqType,
		Status:    models.ServiceStatusPending,
		CreatedAt: time.Now(),
		Details:   details,
	}

	requests = append(requests, newRequest)

	if err := s.serviceRequestRepo.SaveServiceRequests(requests); err != nil {
		return fmt.Errorf("failed to save request: %w", err)
	}

	fmt.Printf(color.GreenString("✅%s request submitted for room %d!\n"), reqType, roomNum)
	return nil
}

func (srs *ServiceRequestService) GetPendingRequestCount() (int, error) {
	requests, err := srs.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, r := range requests {
		if r.Status == models.ServiceStatusPending {
			count++
		}
	}
	return count, nil
}

func (s *ServiceRequestService) GetUnassignedServiceRequest() ([]models.ServiceRequest, error) {
	return s.serviceRequestRepo.GetUnassignedRequests()
}

func (s *ServiceRequestService) CancelServiceRequestByID(reqID string) error {
	req, err := s.serviceRequestRepo.GetServiceRequestByReqID(reqID)
	if err != nil {
		return fmt.Errorf("service request not found: %w", err)
	}

	req.Status = models.ServiceStatusCancelled
	req.UpdatedAt = time.Now()

	return s.serviceRequestRepo.UpdateServiceRequest(req)
}

func (s *ServiceRequestService) UpdateServiceRequestStatus(reqID string, status models.ServiceStatus) error {
	req, err := s.serviceRequestRepo.GetServiceRequestByReqID(reqID)
	if err != nil {
		return fmt.Errorf("service request not found: %v", err)
	}

	req.Status = status
	req.UpdatedAt = time.Now()

	return s.serviceRequestRepo.UpdateServiceRequest(req)
}


func (s *ServiceRequestService) UpdateServiceRequestAssignment(reqID string, isAssigned bool) error {
	return s.serviceRequestRepo.UpdateIsAssigned(reqID, isAssigned)
}

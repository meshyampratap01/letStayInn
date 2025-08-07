package servicerequest

import (
	"context"
	"fmt"
	"time"

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

func NewServiceRequestService(bookingRepo bookingRepository.BookingRepository, serviceRequestRepo serviceRequestRepository.ServiceRequestRepository) *ServiceRequestService {
	return &ServiceRequestService{
		bookingRepo:        bookingRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (s *ServiceRequestService) ServiceRequest(ctx context.Context, roomNum int, reqType models.ServiceType) error {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return fmt.Errorf("failed to load bookings: %w", err)
	}

	hasValidBooking := false
	for _, b := range bookings {
		if b.UserID == ctx.Value(contextkeys.UserIDKey) && b.RoomNum == roomNum &&
			b.Status != models.BookingStatusCancelled {
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
	for _, r := range requests {
		if r.UserID == ctx.Value(contextkeys.UserIDKey) && r.RoomNum == roomNum && r.Type == reqType {
			return fmt.Errorf("you have already requested %s for this room", reqType)
		}
	}

	newRequest := models.ServiceRequest{
		ID:        utils.NewUUID(),
		UserID:    ctx.Value(contextkeys.UserIDKey).(string),
		RoomNum:    roomNum,
		Type:      reqType,
		Status:    models.ServiceStatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	requests = append(requests, newRequest)

	if err := s.serviceRequestRepo.SaveServiceRequests(requests); err != nil {
		return fmt.Errorf("failed to save request: %w", err)
	}

	fmt.Printf("%s request submitted for room %d!\n", reqType, roomNum)
	return nil
}

func (srs *ServiceRequestService) GetPendingRequestCount() (int, error) {
	requests, err := srs.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, r := range requests {
		if r.Status == "pending" {
			count++
		}
	}
	return count, nil
}

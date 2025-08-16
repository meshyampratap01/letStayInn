package employeeService

import (
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
)

type EmployeeService struct {
	userRepo           userRepository.UserRepository
	roomRepo           roomRepository.IRoomRepository
	bookingRepo        bookingRepository.BookingRepository
	serviceRequestRepo serviceRequestRepository.ServiceRequestRepository
}

func NewEmployeeService(userRepo userRepository.UserRepository, roomRepo roomRepository.IRoomRepository, bookingRepo bookingRepository.BookingRepository, serviceRequestRepo serviceRequestRepository.ServiceRequestRepository) IEmployeeService {
	return &EmployeeService{
		userRepo:           userRepo,
		roomRepo:           roomRepo,
		bookingRepo:        bookingRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (es *EmployeeService) GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) {
	return es.serviceRequestRepo.GetAssignedServiceRequests(employeeID)
}

func (es *EmployeeService) UpdateServiceRequestStatus(requestID string, newStatus models.ServiceStatus) error {
    requests, err := es.serviceRequestRepo.LoadServiceRequests()
    if err != nil {
        return err
    }

    var requestToUpdate *models.ServiceRequest
    for i := range requests {
        if requests[i].ID == requestID {
            requestToUpdate = &requests[i]
            break
        }
    }

    if requestToUpdate == nil {
        return fmt.Errorf("service request with ID %s not found", requestID)
    }

    requestToUpdate.Status = newStatus

    return es.serviceRequestRepo.UpdateServiceRequest(requestToUpdate)
}
func (es *EmployeeService) ToggleAvailability(userID string) error {
	return es.userRepo.ToggleStaffAvailability(userID)
}

func (es *EmployeeService) GetAvailability(userID string) (bool, error) {
	return es.userRepo.GetStaffAvailability(userID)
}

func (es *EmployeeService) GetRoomNumberByBookingID(bookingID string) (string, error) {
	booking, err := es.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return "", fmt.Errorf("failed to find booking: %w", err)
	}

	roomNum, err := es.roomRepo.GetRoomNumberByBookingID(booking.ID)
	if err != nil {
		return "", fmt.Errorf("failed to find room: %w", err)
	}

	return roomNum, nil
}

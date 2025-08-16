package employeeService

import "github.com/meshyampratap01/letStayInn/internal/models"

type IEmployeeService interface {
	GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) 
	ToggleAvailability(userID string) error
	GetAvailability(userID string) (bool, error)
	GetRoomNumberByBookingID(string) (string, error)
	UpdateServiceRequestStatus(requestID string, newStatus models.ServiceStatus) error
}

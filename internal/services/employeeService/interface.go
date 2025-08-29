package employeeService

import "github.com/meshyampratap01/letStayInn/internal/models"
//go:generate mockgen -source=interface.go -destination=../../mocks/mock_employeeService.go -package=mocks

type IEmployeeService interface {
	GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) 
	ToggleAvailability(userID string) error
	GetAvailability(userID string) (bool, error)
	GetRoomNumberByBookingID(string) (string, error)
	UpdateServiceRequestStatus(requestID string, newStatus models.ServiceStatus) error
}

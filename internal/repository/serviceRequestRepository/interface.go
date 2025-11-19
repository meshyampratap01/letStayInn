package serviceRequestRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type ServiceRequestRepository interface {
	LoadServiceRequests() ([]models.ServiceRequest, error)
	SaveServiceRequest(req models.ServiceRequest) error
	GetPendingServiceRequests() ([]models.ServiceRequest, error)
	GetServiceRequestByRoomNum(roomNum int) (*models.ServiceRequest, error)
	GetServiceRequestByReqID(id string) (*models.ServiceRequest, error)
	UpdateServiceRequest(req *models.ServiceRequest) error
	GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error)
	DeleteRoomRequests(roomNum int) error
}

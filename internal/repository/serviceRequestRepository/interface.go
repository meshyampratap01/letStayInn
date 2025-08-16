package serviceRequestRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type ServiceRequestRepository interface {
	LoadServiceRequests() ([]models.ServiceRequest, error)
	SaveServiceRequests([]models.ServiceRequest) error
	GetUnassignedRequests() ([]models.ServiceRequest, error)
	GetServiceRequestByRoomNum(roomNum int) (*models.ServiceRequest, error)
	GetServiceRequestByReqID(id string) (*models.ServiceRequest, error)
	UpdateServiceRequest(req *models.ServiceRequest) error
	GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error)
	UpdateIsAssigned(reqID string, isAssigned bool) error
}


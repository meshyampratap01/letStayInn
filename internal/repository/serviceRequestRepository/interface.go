package serviceRequestRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type ServiceRequestRepository interface {
	LoadServiceRequests() ([]models.ServiceRequest, error)
	SaveServiceRequests([]models.ServiceRequest) error
	GetUnassignedRequests() ([]models.ServiceRequest, error)
}

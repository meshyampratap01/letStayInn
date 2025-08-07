package serviceRequestRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type ServiceRequestRepository interface {
	LoadServiceRequests() ([]models.ServiceRequest, error)
	SaveServiceRequests(requests []models.ServiceRequest) error
	GetUnassignedRequests() ([]models.ServiceRequest, error)
}

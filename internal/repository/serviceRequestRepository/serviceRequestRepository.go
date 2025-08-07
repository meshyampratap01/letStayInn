package serviceRequestRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileServiceRequestRepository struct{}

func NewFileServiceRequestRepository() ServiceRequestRepository{
	return &FileServiceRequestRepository{}
}

func (db *FileServiceRequestRepository)LoadServiceRequests() ([]models.ServiceRequest, error) {
	var requests []models.ServiceRequest
	err := storage.ReadJson(config.ServiceRequestFile, &requests)
	return requests, err
}

func (db *FileServiceRequestRepository) SaveServiceRequests(requests []models.ServiceRequest) error {
	return storage.WriteJson(config.ServiceRequestFile, requests)
}

func (r *FileServiceRequestRepository) GetUnassignedRequests() ([]models.ServiceRequest, error) {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return nil, err
	}

	var unassigned []models.ServiceRequest
	for _, req := range requests {
		if !req.IsAssigned {
			unassigned = append(unassigned, req)
		}
	}
	return unassigned, nil
}


package serviceRequestRepository

import (
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileServiceRequestRepository struct{}

func NewFileServiceRequestRepository() ServiceRequestRepository {
	return &FileServiceRequestRepository{}
}

func (db *FileServiceRequestRepository) LoadServiceRequests() ([]models.ServiceRequest, error) {
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

func (r *FileServiceRequestRepository) GetServiceRequestByRoomNum(roomNum int) (*models.ServiceRequest, error) {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return nil, err
	}

	for _, req := range requests {
		if req.RoomNum == roomNum {
			copy := req
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("service request for room %d not found", roomNum)
}

func (r *FileServiceRequestRepository) UpdateServiceRequest(req *models.ServiceRequest) error {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return err
	}

	updated := false
	for i := range requests {
		if requests[i].ID == req.ID {
			requests[i] = *req
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("service request with id %s not found", req.ID)
	}

	return r.SaveServiceRequests(requests)
}

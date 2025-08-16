package serviceRequestRepository

import (
	"fmt"
	"time"

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
		if (!req.IsAssigned || req.Status == models.ServiceStatusPending) && req.Status!=models.ServiceStatusCancelled {
			unassigned = append(unassigned, req)
		}
	}
	return unassigned, nil
}

func (r *FileServiceRequestRepository) UpdateIsAssigned(reqID string, isAssigned bool) error {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return fmt.Errorf("failed to load service requests: %w", err)
	}

	updated := false
	for i, req := range requests {
		if req.ID == reqID {
			requests[i].IsAssigned = isAssigned
			requests[i].UpdatedAt = time.Now()
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("service request with ID %s not found", reqID)
	}

	if err := r.SaveServiceRequests(requests); err != nil {
		return fmt.Errorf("failed to save updated service requests: %w", err)
	}

	return nil
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

func (r *FileServiceRequestRepository) GetServiceRequestByReqID(id string) (*models.ServiceRequest, error) {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return nil, err
	}
	for _, req := range requests {
		if req.ID == id {
			return &req, nil
		}
	}
	return nil, fmt.Errorf("service request with ID %s not found", id)
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

func (r *FileServiceRequestRepository) GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) {
	requests, err := r.LoadServiceRequests()
	if err != nil {
		return nil, err
	}
	var assigned []models.ServiceRequest
	for _, req := range requests {
		if req.AssignedTo == employeeID {
			assigned = append(assigned, req)
		}
	}
	return assigned, nil
}

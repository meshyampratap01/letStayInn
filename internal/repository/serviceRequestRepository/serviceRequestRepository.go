package serviceRequestRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/gorm"
)

type GormServiceRequestRepository struct {
	db *gorm.DB
}

func NewGormServiceRequestRepository(db *gorm.DB) ServiceRequestRepository {
	return &GormServiceRequestRepository{db: db}
}

func (r *GormServiceRequestRepository) LoadServiceRequests() ([]models.ServiceRequest, error) {
	var requests []models.ServiceRequest
	err := r.db.Find(&requests).Error
	return requests, err
}

func (r *GormServiceRequestRepository) SaveServiceRequest(req models.ServiceRequest) error {
	return r.db.Create(&req).Error
}

func (r *GormServiceRequestRepository) GetUnassignedRequests() ([]models.ServiceRequest, error) {
	var requests []models.ServiceRequest
	err := r.db.Where("(is_assigned = ? OR status = ?) AND status != ? AND status != ?", false, models.ServiceStatusPending, models.ServiceStatusCancelled, models.ServiceStatusDone).Find(&requests).Error
	return requests, err
}

func (r *GormServiceRequestRepository) GetServiceRequestByRoomNum(roomNum int) (*models.ServiceRequest, error) {
	var req models.ServiceRequest
	err := r.db.Where("room_num = ?", roomNum).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *GormServiceRequestRepository) GetServiceRequestByReqID(id string) (*models.ServiceRequest, error) {
	var req models.ServiceRequest
	err := r.db.First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *GormServiceRequestRepository) UpdateServiceRequest(req *models.ServiceRequest) error {
	return r.db.Save(req).Error
}

func (r *GormServiceRequestRepository) GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) {
	var requests []models.ServiceRequest
	err := r.db.Where("assigned_to = ?", employeeID).Find(&requests).Error
	return requests, err
}

func (r *GormServiceRequestRepository) UpdateIsAssigned(reqID string, isAssigned bool) error {
	return r.db.Model(&models.ServiceRequest{}).Where("id = ?", reqID).Update("is_assigned", isAssigned).Error
}

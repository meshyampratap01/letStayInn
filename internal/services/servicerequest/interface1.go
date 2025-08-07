package servicerequest

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

// ServiceRequestService defines behavior for handling service requests.
type IServiceRequestService interface {
	ServiceRequest(ctx context.Context, roomNum int, reqType models.ServiceType) error
	GetPendingRequestCount() (int, error)
	ViewAllServiceRequests() ([]models.ServiceRequest, error)
	ViewUnassignedServiceRequest() ([]models.ServiceRequest,error)
}

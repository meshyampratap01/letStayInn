package servicerequest

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

type IServiceRequestService interface {
	ServiceRequestGetter(context.Context,int,models.ServiceType,string) error
	GetPendingRequestCount() (int, error)
	GetUnassignedServiceRequest() ([]models.ServiceRequest,error)
	UpdateServiceRequestStatus(reqID string, status models.ServiceStatus) error
	CancelServiceRequestByID(reqID string) error
	UpdateServiceRequestAssignment(reqID string, isAssigned bool) error
}

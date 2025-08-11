package servicerequest

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

type IServiceRequestService interface {
	ServiceRequestGetter(context.Context,int,models.ServiceType,string) error
	GetPendingRequestCount() (int, error)
	ViewAllServiceRequests() ([]models.ServiceRequest, error)
	ViewUnassignedServiceRequest() ([]models.ServiceRequest,error)
	CancelServiceRequestByRoomNum(roomNum int) error
}

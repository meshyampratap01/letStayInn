package managerservice

import "github.com/meshyampratap01/letStayInn/internal/models"

type IManagerService interface {
	UpdateEmployeeAvailability(email string, available bool) error
	GetTotalEmployees() (int, error)
	GetAllEmployees() ([]models.User, error)
	DeleteEmployeeByEmail(email string) error
	GetAvailableStaffByTaskType(string) ([]models.User, error)
	AssignServiceRequest(reqID string, empID string) error
	PrintHotelReport() error
	ViewAllFeedback() ([]models.Feedback, error)
}

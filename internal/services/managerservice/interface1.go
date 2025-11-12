package managerservice

import "github.com/meshyampratap01/letStayInn/internal/models"

//go:generate mockgen -source=interface1.go -destination=../../mocks/mock_managerService.go -package=mocks

type IManagerService interface {
	UpdateEmployeeAvailability(email string, available bool) error
	GetTotalEmployees() (int, error)
	GetAllEmployees() ([]models.User, error)
	DeleteEmployeeByID(id string) error
	GetAvailableStaffByTaskType(string) ([]models.User, error)
	AssignServiceRequest(reqID string, empID string) error
	ViewAllFeedback() ([]models.Feedback, error)
}

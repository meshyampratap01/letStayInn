package employeeService

import "github.com/meshyampratap01/letStayInn/internal/models"

type IEmployeeService interface {
	ViewAssignedTasks(string) ([]models.Task,error)
	UpdateTaskStatus(string,models.TaskStatus) error
	ToggleAvailability(userID string) error
	GetAvailability(userID string) (bool, error)
	GetRoomNumberByBookingID(string) (string, error)
}
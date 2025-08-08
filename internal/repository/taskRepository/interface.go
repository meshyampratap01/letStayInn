package taskRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type TaskRepository interface {
	SaveAllTasks([]models.Task) error
	SaveTask(models.Task) error
	GetAllTask() ([]models.Task,error)
	GetTasksByStaffID(string) ([]models.Task, error)
	UpdateTaskStatus(taskID string, status models.TaskStatus) error
}

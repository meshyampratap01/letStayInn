package taskRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type TaskRepository interface {
	SaveTask(models.Task) error
}

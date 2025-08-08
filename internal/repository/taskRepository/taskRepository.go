package taskRepository

import (
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileTaskRepository struct{}

func NewFileTaskRepository() TaskRepository {
	return &FileTaskRepository{}
}

func (tr *FileTaskRepository) SaveTask(task models.Task) error {
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile, tasks)
	if err != nil {
		return err
	}

	tasks = append(tasks, task)

	return storage.WriteJson(config.TasksFile, tasks)
}

func (tr *FileTaskRepository) GetAllTask() ([]models.Task, error) {
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *FileTaskRepository) SaveAllTasks(tasks []models.Task) error {
	return storage.WriteJson(config.TasksFile, tasks)
}

func (repo *FileTaskRepository) GetTasksByStaffID(staffID string) ([]models.Task, error) {
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile, &tasks)
	if err != nil {
		return nil, err
	}

	var assigned []models.Task
	for _, t := range tasks {
		if t.AssignedTo == staffID {
			assigned = append(assigned, t)
		}
	}
	return assigned, nil
}

func (repo *FileTaskRepository) UpdateTaskStatus(taskID string, status models.TaskStatus) error {
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile, &tasks)
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Status = status
			return storage.WriteJson(config.TasksFile, tasks)
		}
	}
	return fmt.Errorf("task not found")
}

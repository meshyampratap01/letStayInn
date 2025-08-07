package taskRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileTaskRepository struct{}

func NewFileTaskRepository() TaskRepository {
	return &FileTaskRepository{}
}

func (tr *FileTaskRepository) SaveTask(task models.Task) error{
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile,tasks)
	if err!=nil{
		return err
	}

	tasks = append(tasks,task)

	return storage.WriteJson(config.TasksFile,tasks)
}

func (tr *FileTaskRepository) GetAllTask() ([]models.Task,error){
	var tasks []models.Task
	err := storage.ReadJson(config.TasksFile,tasks)
	if err!=nil{
		return nil,err
	}
	return tasks,nil
}

func (tr *FileTaskRepository) SaveAllTasks(tasks []models.Task) error{
	return storage.WriteJson(config.TasksFile,tasks)
}
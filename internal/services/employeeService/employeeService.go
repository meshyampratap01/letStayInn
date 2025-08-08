package employeeService

import (
	"fmt"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/taskRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
)

type EmployeeService struct {
	taskRepo           taskRepository.TaskRepository
	userRepo           userRepository.UserRepository
	roomRepo           roomRepository.IRoomRepository
	bookingRepo        bookingRepository.BookingRepository
	serviceRequestRepo serviceRequestRepository.ServiceRequestRepository
}

func NewEmployeeService(taskRepo taskRepository.TaskRepository, userRepo userRepository.UserRepository, roomRepo roomRepository.IRoomRepository, bookingRepo bookingRepository.BookingRepository, serviceRequestRepo serviceRequestRepository.ServiceRequestRepository) IEmployeeService {
	return &EmployeeService{
		taskRepo:           taskRepo,
		userRepo:           userRepo,
		roomRepo:           roomRepo,
		bookingRepo:        bookingRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (es *EmployeeService) ViewAssignedTasks(userID string) ([]models.Task, error) {
	return es.taskRepo.GetTasksByStaffID(userID)
}

func (s *EmployeeService) UpdateTaskStatus(taskID string, newStatus models.TaskStatus) error {
	// 1. Load all tasks
	tasks, err := s.taskRepo.GetAllTask()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// 2. Find the task
	var updatedTask *models.Task
	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now()
			updatedTask = &tasks[i]
			break
		}
	}
	if updatedTask == nil {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// 3. Save updated tasks list
	if err := s.taskRepo.SaveAllTasks(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	// 4. Also update the related ServiceRequest
	requests, err := s.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return fmt.Errorf("failed to load service requests: %w", err)
	}

	for i := range requests {
		if requests[i].ID == updatedTask.RequestID {
			// Map TaskStatus → ServiceStatus
			switch newStatus {
			case models.TaskStatusPending:
				requests[i].Status = models.ServiceStatusPending
			case models.TaskStatusInProgress:
				requests[i].Status = models.ServiceStatusInProgrss
			case models.TaskStatusDone:
				requests[i].Status = models.ServiceStatusDone
			}
			break
		}
	}

	// 5. Save updated requests
	if err := s.serviceRequestRepo.SaveServiceRequests(requests); err != nil {
		return fmt.Errorf("failed to save service requests: %w", err)
	}

	return nil
}

func (es *EmployeeService) ToggleAvailability(userID string) error {
	return es.userRepo.ToggleStaffAvailability(userID)
}

func (es *EmployeeService) GetAvailability(userID string) (bool, error) {
	return es.userRepo.GetStaffAvailability(userID)
}

func (es *EmployeeService) GetRoomNumberByBookingID(bookingID string) (string, error) {
	booking, err := es.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return "", fmt.Errorf("failed to find booking: %w", err)
	}

	roomNum, err := es.roomRepo.GetRoomNumberByBookingID(booking.ID)
	if err != nil {
		return "", fmt.Errorf("failed to find room: %w", err)
	}

	return roomNum, nil
}

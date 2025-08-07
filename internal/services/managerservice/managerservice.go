package managerservice

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/taskRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type ManagerService struct {
	userRepo userRepository.UserRepository
	taskRepo taskRepository.TaskRepository
	serviceRequestRepo 	serviceRequestRepository.ServiceRequestRepository
}

func NewManagerService(userRepo userRepository.UserRepository,taskRepo taskRepository.TaskRepository,serviceRequestRepo serviceRequestRepository.ServiceRequestRepository) IManagerService {
	return &ManagerService{
		userRepo: userRepo,
		taskRepo: taskRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (ms *ManagerService) UpdateEmployeeAvailability(email string, available bool) error {
	users, err := ms.userRepo.GetAllUsers()
	if err != nil {
		return err
	}

	found := false
	for i, u := range users {
		if u.Email == email {
			users[i].Available = available
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("employee with email %s not found", email)
	}

	err = ms.userRepo.SaveAllUsers(users)
	if err != nil {
		return err
	}
	return nil
}

func (ms *ManagerService) GetAllEmployees() ([]models.User, error) {
	users, err := ms.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var employees []models.User
	for _, user := range users {
		if user.Role == models.RoleKitchenStaff || user.Role == models.RoleCleaningStaff {
			employees = append(employees, user)
		}
	}
	return employees, nil
}

func (ms *ManagerService) GetTotalEmployees() (int, error) {
	users, err := ms.userRepo.GetAllUsers()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, u := range users {
		if u.Role == models.RoleCleaningStaff || u.Role == models.RoleKitchenStaff {
			count++
		}
	}
	return count, nil
}

func (ms *ManagerService) DeleteEmployeeByEmail(email string) error {
	users, err := ms.userRepo.GetAllUsers()
	if err != nil {
		return err
	}

	found := false
	var updatedUsers []models.User
	for _, user := range users {
		if user.Email == email {
			if user.Role != models.RoleKitchenStaff && user.Role != models.RoleCleaningStaff {
				return errors.New("cannot delete non-employee user")
			}
			found = true
			continue // skipping the employee to delete
		}
		updatedUsers = append(updatedUsers, user)
	}

	if !found {
		return errors.New("employee not found")
	}

	return ms.userRepo.SaveAllUsers(updatedUsers)
}

func (ms *ManagerService) GetAvailableStaffByTaskType(taskType string) ([]models.User, error) {
	var role models.Role
	switch taskType {
	case string(models.ServiceTypeCleaning):
		role = models.RoleCleaningStaff
	case string(models.ServiceTypeFood):
		role = models.RoleKitchenStaff
	default:
		return nil, errors.New("invalid task type")
	}

	allUsers, err := ms.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var availableStaff []models.User
	for _, user := range allUsers {
		if user.Role == role && user.Available {
			availableStaff = append(availableStaff, user)
		}
	}

	return availableStaff, nil
}

func (s *ManagerService) AssignTask(taskType, bookingID, details, staffID string) error {
	var tType models.TaskType
	switch strings.ToLower(taskType) {
	case string(models.ServiceTypeCleaning):
		tType = models.TaskTypeCleaning
	case string(models.ServiceTypeFood):
		tType = models.TaskTypeFood
	default:
		return errors.New("invalid task type")
	}

	newTask := models.Task{
		ID:         utils.NewUUID(),
		Type:       tType,
		AssignedTo: staffID,
		BookingID:  bookingID,
		Details:    details,
		Status:     models.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.taskRepo.SaveTask(newTask)
}

func (s *ManagerService) ViewUnassignedServiceRequest() ([]models.ServiceRequest,error){
	return s.serviceRequestRepo.GetUnassignedRequests()
}

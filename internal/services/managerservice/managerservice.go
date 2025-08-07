package managerservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/taskRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type ManagerService struct {
	userRepo userRepository.UserRepository
	taskRepo taskRepository.TaskRepository
	serviceRequestRepo 	serviceRequestRepository.ServiceRequestRepository
	roomRepo		roomRepository.IRoomRepository
	bookingRepo		bookingRepository.BookingRepository
}

func NewManagerService(userRepo userRepository.UserRepository,taskRepo taskRepository.TaskRepository,serviceRequestRepo serviceRequestRepository.ServiceRequestRepository,roomRepo roomRepository.IRoomRepository,bookingRepo bookingRepository.BookingRepository) IManagerService {
	return &ManagerService{
		userRepo: userRepo,
		taskRepo: taskRepo,
		serviceRequestRepo: serviceRequestRepo,
		roomRepo: roomRepo,
		bookingRepo: bookingRepo,
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

func (s *ManagerService) AssignTaskFromServiceRequest(
	requestID, bookingID, details, staffID string,
) error {
	// 1. Fetch all service requests
	requests, err := s.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return fmt.Errorf("failed to load service requests: %w", err)
	}

	// 2. Find the request by ID
	var targetRequest *models.ServiceRequest
	for i := range requests {
		if requests[i].ID == requestID {
			targetRequest = &requests[i]
			break
		}
	}
	if targetRequest == nil {
		return fmt.Errorf("service request with ID %s not found", requestID)
	}
	if targetRequest.IsAssigned {
		return fmt.Errorf("service request already assigned")
	}

	// 3. Create the new task
	task := models.Task{
		ID:         utils.NewUUID(),
		BookingID:  bookingID,
		AssignedTo: staffID,
		Type:       models.TaskType(targetRequest.Type),
		Details:    details,
		Status:     models.TaskStatusPending,
		CreatedAt:  time.Now(),
	}

	// 4. Save the new task
	tasks, err := s.taskRepo.GetAllTask()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	tasks = append(tasks, task)

	if err := s.taskRepo.SaveAllTasks(tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// 5. Mark the service request as assigned
	targetRequest.IsAssigned = true
	if err := s.serviceRequestRepo.SaveServiceRequests(requests); err != nil {
		return fmt.Errorf("failed to update service requests: %w", err)
	}

	return nil
}



func (s *ManagerService) PrintHotelReport() error {
	// Fetch data from repos
	rooms, _ := s.roomRepo.GetAllRooms()
	availableRooms, _ := s.roomRepo.GetAvailableRooms()
	guests, _ := s.userRepo.GetAllUsers()
	employee, _ := s.GetAllEmployees()
	bookings, _ := s.bookingRepo.GetAllBookings()
	serviceRequests, _ := s.serviceRequestRepo.LoadServiceRequests()
	tasks, _ := s.taskRepo.GetAllTask()


	// Count unassigned service requests
	unassignedRequests := 0
	for _, req := range serviceRequests {
		if !req.IsAssigned {
			unassignedRequests++
		}
	}

	// Task Summary
	taskSummary := make(map[models.TaskType]map[models.TaskStatus]int)
	for _, task := range tasks {
		if _, exists := taskSummary[task.Type]; !exists {
			taskSummary[task.Type] = make(map[models.TaskStatus]int)
		}
		taskSummary[task.Type][task.Status]++
	}

	// Print to CLI
	fmt.Println("\n--- Hotel Report ---")
	fmt.Printf("Total Rooms: %d\n", len(rooms))
	fmt.Printf("Available Rooms: %d\n", len(availableRooms))
	fmt.Printf("Total Guests: %d\n", len(guests))
	fmt.Printf("Total Staff: %d\n", len(employee))
	fmt.Printf("Total Bookings: %d\n", len(bookings))
	fmt.Printf("Unassigned Service Requests: %d\n", unassignedRequests)

	fmt.Println("\n--- Task Summary ---")
	for taskType, statuses := range taskSummary {
		fmt.Printf("%s:\n", taskType)
		for status, count := range statuses {
			fmt.Printf("  %s: %d\n", status, count)
		}
	}

	return nil
}

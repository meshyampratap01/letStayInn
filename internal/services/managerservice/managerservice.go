package managerservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
)

type ManagerService struct {
	userRepo           userRepository.UserRepository
	serviceRequestRepo serviceRequestRepository.ServiceRequestRepository
	roomRepo           roomRepository.IRoomRepository
	bookingRepo        bookingRepository.BookingRepository
	feedbackRepo       feedbackRepository.FeedbackRepository
}

func NewManagerService(userRepo userRepository.UserRepository, serviceRequestRepo serviceRequestRepository.ServiceRequestRepository, roomRepo roomRepository.IRoomRepository, bookingRepo bookingRepository.BookingRepository, feedbackRepo feedbackRepository.FeedbackRepository) IManagerService {
	return &ManagerService{
		userRepo:           userRepo,
		serviceRequestRepo: serviceRequestRepo,
		roomRepo:           roomRepo,
		bookingRepo:        bookingRepo,
		feedbackRepo:       feedbackRepo,
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
				return errors.New("user with this email is not an employee")
			}
			found = true
			continue
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

func (s *ManagerService) AssignServiceRequest(reqID string, empID string) error {
	req, err := s.serviceRequestRepo.GetServiceRequestByReqID(reqID)
	if err != nil {
		return err
	}
	req.AssignedTo = empID
	req.IsAssigned = true
	req.UpdatedAt = time.Now()
	return s.serviceRequestRepo.UpdateServiceRequest(req)
}

func (s *ManagerService) PrintHotelReport() error {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return fmt.Errorf("error fetching rooms: %v", err)
	}

	availableRooms, err := s.roomRepo.GetAvailableRooms()
	if err != nil {
		return fmt.Errorf("error fetching available rooms: %v", err)
	}

	employees, err := s.GetAllEmployees()
	if err != nil {
		return fmt.Errorf("error fetching employees: %v", err)
	}

	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return fmt.Errorf("error fetching bookings: %v", err)
	}

	serviceRequests, err := s.serviceRequestRepo.LoadServiceRequests()
	if err != nil {
		return fmt.Errorf("error fetching service requests: %v", err)
	}

	unassignedRequests := 0
	for _, req := range serviceRequests {
		if !req.IsAssigned {
			unassignedRequests++
		}
	}


	requestStatusSummary := make(map[models.ServiceStatus]int)
	for _, req := range serviceRequests {
		requestStatusSummary[req.Status]++
	}

	fmt.Println("\n--- Hotel Report ---")
	fmt.Printf("Total Rooms: %d\n", len(rooms))
	fmt.Printf("Available Rooms: %d\n", len(availableRooms))
	fmt.Printf("Total Staff: %d\n", len(employees))
	fmt.Printf("Total Bookings: %d\n", len(bookings))
	fmt.Printf("Unassigned Service Requests: %d\n", unassignedRequests)

	fmt.Println("\n--- Service Request Summary ---")
	if len(requestStatusSummary) == 0 {
		fmt.Println("No Service Requests to show.")
	} else {
		for status, count := range requestStatusSummary {
			fmt.Printf("%s: %d\n", status, count)
		}
	}

	return nil
}


func (ms *ManagerService) ViewAllFeedback() ([]models.Feedback, error) {
	return ms.feedbackRepo.GetAllFeedback()
}

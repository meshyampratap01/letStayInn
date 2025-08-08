package handlers

import (
	"context"
	"fmt"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

type DashboardHandler struct {
	RoomService           roomService.RoomServiceManager
	BookingService        bookingService.BookingManager
	FeedbackService       feedbackService.FeedbackServiceManager
	ServiceRequestService servicerequest.IServiceRequestService
	BookingHandler        *BookingHandler
	ServiceRequestHandler *ServiceRequestHandler
	managerHandler        *ManagerHandler
	employeeService       employeeService.IEmployeeService
	employeeHandler       *EmployeeHandler
}

func NewDashboardHandler(
	roomSvc roomService.RoomServiceManager,
	bookingSvc bookingService.BookingManager,
	feedbackSvc feedbackService.FeedbackServiceManager,
	serviceReqSvc servicerequest.IServiceRequestService,
	bh *BookingHandler,
	ServiceRequestHandler *ServiceRequestHandler,
	managerHandler *ManagerHandler,
	employeeService employeeService.IEmployeeService,
	employeeHandler *EmployeeHandler,
) *DashboardHandler {
	return &DashboardHandler{
		RoomService:           roomSvc,
		BookingService:        bookingSvc,
		FeedbackService:       feedbackSvc,
		ServiceRequestService: serviceReqSvc,
		BookingHandler:        bh,
		ServiceRequestHandler: ServiceRequestHandler,
		managerHandler:        managerHandler,
		employeeService:       employeeService,
		employeeHandler:	   employeeHandler,
	}
}

func (h *DashboardHandler) LoadDashboard(ctx context.Context) {
	switch ctx.Value(contextkeys.UserRoleKey) {
	case models.RoleGuest:
		h.guestDashboard(ctx)
	case models.RoleKitchenStaff:
		h.EmployeeDashboard(ctx)
	case models.RoleCleaningStaff:
		h.EmployeeDashboard(ctx)
	case models.RoleManager:
		h.managerDashboard(ctx)
	default:
		fmt.Println("Unknown role.")
	}
}

func (h *DashboardHandler) guestDashboard(ctx context.Context) {
	for {
		fmt.Println("\n--- Guest Dashboard ---")
		fmt.Println("1. View Available Rooms")
		fmt.Println("2. Book Room")
		fmt.Println("3. Cancel Booking")
		fmt.Println("4. View My Bookings")
		fmt.Println("5. Request Food")
		fmt.Println("6. Request Room Cleaning")
		fmt.Println("7. Give Feedback")
		fmt.Println("8. Logout")
		fmt.Print("Select option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			h.BookingHandler.ViewRoomsHandler()
		case 2:
			h.BookingHandler.BookRoomHandler(ctx)
		case 3:
			h.BookingHandler.CancelBookingHandler(ctx)
		case 4:
			h.BookingHandler.ViewMyBookingsHandler(ctx)
		case 5:
			h.ServiceRequestHandler.ServiceRequestHandler(ctx, models.ServiceTypeFood)
		case 6:
			h.ServiceRequestHandler.ServiceRequestHandler(ctx, models.ServiceTypeCleaning)
		case 7:
			err := h.FeedbackService.SubmitFeedback(ctx)
			if err != nil {
				fmt.Println("Error submitting feedback:", err)
			}
		case 8:
			fmt.Println("Logging out...")
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func (h *DashboardHandler) managerDashboard(ctx context.Context) {
	for {
		fmt.Println("\n--- Manager Dashboard ---")
		fmt.Println("1. View Dashboard Summary")
		fmt.Println("2. Room Management")
		fmt.Println("3. View Bookings and Guests")
		fmt.Println("4. Manage Staff")
		fmt.Println("5. View All Guest Service Requests")
		fmt.Println("6. View Unassigned Guest Service Requests")
		fmt.Println("7. Assign Task to Employee")
		fmt.Println("8. Generate Reports")
		fmt.Println("9. View Guest Feedback")
		fmt.Println("10. Logout")

		fmt.Print("Select option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			h.managerHandler.ManagerDashboardSummary()

		case 2:
		RoomMgmtLoop:
			for {
				fmt.Println("\n--- Room Management ---")
				fmt.Println("1. List Rooms")
				fmt.Println("2. Add Room")
				fmt.Println("3. Update Room")
				fmt.Println("4. Delete Room")
				fmt.Println("5. Back")
				fmt.Print("Select option: ")
				var rchoice int
				fmt.Scanln(&rchoice)
				switch rchoice {
				case 1:
					h.managerHandler.ListRooms()
				case 2:
					h.managerHandler.AddRoom()
				case 3:
					h.managerHandler.UpdateRoom()
				case 4:
					h.managerHandler.DeleteRoom()
				case 5:
					break RoomMgmtLoop
				default:
					fmt.Println("Invalid option.")
				}
			}

		case 3:
			h.managerHandler.ListBookingsAndGuests()

		case 4:
		EmpMgmtLoop:
			for {
				fmt.Println("\n--- Employee Management ---")
				fmt.Println("1. List Employee")
				fmt.Println("2. Update Employee Availability")
				fmt.Println("3. Delete employee")
				fmt.Println("4. Back")
				fmt.Print("Select option: ")
				var echoice int
				fmt.Scanln(&echoice)
				switch echoice {
				case 1:
					h.managerHandler.ListEmployee()
				case 2:
					h.managerHandler.UpdateEmployeeAvailability()
				case 3:
					h.managerHandler.DeleteEmployee()
				case 4:
					break EmpMgmtLoop
				default:
					fmt.Println("Invalid option.")
				}
			}

		case 5:
			h.managerHandler.ViewAllServiceRequests()
		case 6:
			h.managerHandler.ViewUnassignedServiceRequests()
		case 7:
			h.managerHandler.AssignTasksToEmployees()
		case 8:
			h.managerHandler.GenerateReport()
		case 9:
			h.managerHandler.ViewFeedback(ctx) 
		case 10:
			fmt.Println("Logging out...")
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}


func (h *DashboardHandler) EmployeeDashboard(ctx context.Context) {

	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		fmt.Println("invalid or missing user ID in context")
		return
	}
	for {
		fmt.Println("\n--- Employee Dashboard ---")
		fmt.Println("1. View Assigned Tasks")
		fmt.Println("2. Update Task Status")
		fmt.Println("3. Toggle Availability")
		fmt.Println("4. Exit Dashboard")
		fmt.Print("Select an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			h.employeeHandler.ViewAssignedTasks(userID)
		case 2:
			h.employeeHandler.UpdateTaskStatus(userID)
		case 3:
			h.employeeHandler.ToggleAvailability(userID)
		case 4:
			fmt.Println("Exiting Employee Dashboard...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}


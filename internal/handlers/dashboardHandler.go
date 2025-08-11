package handlers

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

var (
	titleStyle   = color.New(color.Bold, color.FgHiWhite).SprintFunc()
	optionStyle  = color.New(color.FgHiCyan).SprintFunc()
	promptStyle  = color.New(color.FgHiYellow).SprintFunc()
	errStyle     = color.New(color.FgHiRed).SprintFunc()
	successStyle = color.New(color.FgHiGreen).SprintFunc()
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
		employeeHandler:       employeeHandler,
	}
}

func (h *DashboardHandler) LoadDashboard(ctx context.Context) {
	switch ctx.Value(contextkeys.UserRoleKey) {
	case models.RoleGuest:
		h.guestDashboard(ctx)
	case models.RoleKitchenStaff, models.RoleCleaningStaff:
		h.EmployeeDashboard(ctx)
	case models.RoleManager:
		h.managerDashboard(ctx)
	default:
		fmt.Println(errStyle("Unknown role."))
	}
}

func (h *DashboardHandler) guestDashboard(ctx context.Context) {
	for {
		fmt.Println(titleStyle(config.GuestDashboardTitle))
		fmt.Println(optionStyle("\n1.") + " View Available Rooms")
		fmt.Println(optionStyle("2.") + " Book Room")
		fmt.Println(optionStyle("3.") + " Cancel Booking")
		fmt.Println(optionStyle("4.") + " View My Bookings")
		fmt.Println(optionStyle("5.") + " Request Food")
		fmt.Println(optionStyle("6.") + " Request Room Cleaning")
		fmt.Println(optionStyle("7.") + " Give Feedback")
		fmt.Println(optionStyle("8.") + " Logout")
		fmt.Print(promptStyle(config.SelectOption))

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
			if err := h.FeedbackService.SubmitFeedback(ctx); err != nil {
				fmt.Println(errStyle("Error submitting feedback:"), err)
			}
		case 8:
			fmt.Println(successStyle("Logging out..."))
			return
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (h *DashboardHandler) managerDashboard(ctx context.Context) {
	for {
		fmt.Println(titleStyle(config.ManagerDashboardTitle))
		fmt.Println(optionStyle("\n1.") + " View Dashboard Summary")
		fmt.Println(optionStyle("2.") + " Room Management")
		fmt.Println(optionStyle("3.") + " View Bookings and Guests")
		fmt.Println(optionStyle("4.") + " Manage Staff")
		fmt.Println(optionStyle("5.") + " Manage Service Requests")
		fmt.Println(optionStyle("6.") + " Generate Reports")
		fmt.Println(optionStyle("7.") + " View Guest Feedback")
		fmt.Println(optionStyle("8.") + " Logout")
		fmt.Print(promptStyle(config.SelectOption))

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			h.managerHandler.ManagerDashboardSummary()

		case 2:
			h.managerHandler.roomManagementMenu()

		case 3:
			h.managerHandler.ListBookingsAndGuests()

		case 4:
			h.managerHandler.employeeManagementMenu()

		case 5:
			h.managerHandler.serviceRequestManagementMenu()

		case 6:
			h.managerHandler.GenerateReport()

		case 7:
			h.managerHandler.ViewFeedback(ctx)

		case 8:
			fmt.Println(successStyle("Logging out..."))
			return

		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (h *DashboardHandler) EmployeeDashboard(ctx context.Context) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		fmt.Println(errStyle("Invalid or missing user ID in context"))
		return
	}

	for {
		fmt.Println(titleStyle(config.EmployeeDashboardTitle))
		fmt.Println(optionStyle("\n1.") + " View Assigned Tasks")
		fmt.Println(optionStyle("2.") + " Update Task Status")
		fmt.Println(optionStyle("3.") + " Toggle Availability")
		fmt.Println(optionStyle("4.") + " Exit Dashboard")
		fmt.Print(promptStyle(config.SelectOption))

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
			fmt.Println(successStyle("Exiting Employee Dashboard..."))
			return
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}
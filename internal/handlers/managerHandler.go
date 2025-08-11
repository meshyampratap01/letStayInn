package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	managerservice "github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

type ManagerHandler struct {
	roomService           roomService.RoomServiceManager
	bookingService        bookingService.BookingManager
	userService           userService.UserManager
	serviceRequestService servicerequest.IServiceRequestService
	managerService        managerservice.IManagerService
}

func NewManagerHandler(rs roomService.RoomServiceManager, bs bookingService.BookingManager, us userService.UserManager, srs servicerequest.IServiceRequestService, ms managerservice.IManagerService) *ManagerHandler {
	return &ManagerHandler{
		roomService:           rs,
		bookingService:        bs,
		userService:           us,
		serviceRequestService: srs,
		managerService:        ms,
	}
}

func (mh ManagerHandler) ManagerDashboardSummary() {
	totalRooms, _ := mh.roomService.GetTotalRooms()
	availableRooms, _ := mh.roomService.GetTotalAvailableRooms()
	bookedRooms := totalRooms - availableRooms

	totalGuests, _ := mh.userService.GetTotalGuests()
	totalEmp, _ := mh.managerService.GetTotalEmployees()

	pendingRequests, _ := mh.serviceRequestService.GetPendingRequestCount()

	color.Cyan("\n--- Dashboard Summary ---")
	fmt.Printf("Total Rooms: %d\n", totalRooms)
	fmt.Printf("Booked Rooms: %d\n", bookedRooms)
	fmt.Printf("Available Rooms: %d\n", availableRooms)
	fmt.Printf("Total Guests(including those who checked out): %d\n", totalGuests)
	fmt.Printf("Total Employee: %d\n", totalEmp)
	fmt.Printf("Pending Service Requests: %d\n", pendingRequests)
}

func (h *ManagerHandler) ListRooms() {
	rooms, err := h.roomService.GetAllRooms()
	if err != nil {
		color.Red("Error fetching rooms: %v", err)
		return
	}

	if len(rooms) == 0 {
		color.Yellow("No rooms found.")
		return
	}

	color.Cyan("\n--- Room List ---")
	for _, room := range rooms {
		availability := color.RedString("No")
		if room.IsAvailable {
			availability = color.GreenString("Yes")
		}
		fmt.Printf("Number: %d | Type: %s | Price: %.2f | Available: %s | Description: %s\n",
			room.Number, room.Type, room.Price, availability, room.Description)
	}
}

func (h *ManagerHandler) AddRoom() {
	reader := bufio.NewReader(os.Stdin)

	color.Yellow("Enter room number: ")
	var number int
	fmt.Scanln(&number)

	color.Yellow("Enter room type (e.g., Deluxe, Standard): ")
	roomType, _ := reader.ReadString('\n')
	roomType = strings.TrimSpace(roomType)

	color.Yellow("Enter room price: ")
	var price float64
	fmt.Scanln(&price)

	color.Yellow("Is the room available? (yes/no): ")
	availableInput, _ := reader.ReadString('\n')
	availableInput = strings.TrimSpace(strings.ToLower(availableInput))
	isAvailable := availableInput == "yes"

	color.Yellow("Enter room description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	err := h.roomService.AddRoom(number, roomType, price, isAvailable, description)
	if err != nil {
		color.Red("Error adding room: %v", err)
		return
	}

	color.Green("Room added successfully.")
}

func (h *ManagerHandler) UpdateRoom() {
	reader := bufio.NewReader(os.Stdin)

	h.ListRooms()

	fmt.Print(color.YellowString("\nEnter Room Number to update: "))
	var number int
	fmt.Scanln(&number)

	fmt.Print(color.CyanString("\nSelect what you want to update:"))
	fmt.Println("1. Room Type")
	fmt.Println("2. Price")
	fmt.Println("3. Availability")
	fmt.Println("4. Description")
	color.Yellow("Enter choice: ")
	var choice int
	fmt.Scanln(&choice)

	var roomType, description string
	var price float64
	var isAvailable bool

	switch choice {
	case 1:
		fmt.Print(color.YellowString("Enter new Room Type: "))
		roomType, _ = reader.ReadString('\n')
		roomType = strings.TrimSpace(roomType)
	case 2:
		fmt.Print(color.YellowString("Enter new Price: "))
		fmt.Scanln(&price)
	case 3:
		fmt.Print(color.YellowString("Is room available? (yes/no): "))
		availableInput, _ := reader.ReadString('\n')
		isAvailable = strings.TrimSpace(strings.ToLower(availableInput)) == "yes"
	case 4:
		fmt.Print(color.YellowString("Enter new Description: "))
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
	default:
		color.Red("Invalid choice.")
		return
	}

	err := h.roomService.UpdateRoom(number, choice, roomType, price, isAvailable, description)
	if err != nil {
		color.Red("Error updating room: %v", err)
		return
	}

	color.Green("Room updated successfully.")
}

func (h *ManagerHandler) DeleteRoom() {
	h.ListRooms()
	fmt.Print(color.YellowString("\nEnter Room Number to delete: "))
	var number int
	fmt.Scanln(&number)

	err := h.roomService.DeleteRoom(number)
	if err != nil {
		color.Red("Error deleting room: %v", err)
		return
	}

	color.Green("Room deleted successfully.")
}

func (h *ManagerHandler) ListBookingsAndGuests() {
	bookings, err := h.bookingService.GetAllBookingsWithGuests()
	if err != nil {
		color.Red("Error fetching bookings: %v", err)
		return
	}

	color.Cyan("\n--- Bookings and Guests ---")
	if len(bookings) == 0 {
		color.Yellow("No bookings found.")
		return
	}

	for _, b := range bookings {
		fmt.Println("---------------")
		fmt.Printf("Booking ID   : %s\n", b.ID)
		fmt.Printf("Guest Name   : %s\n", b.GuestName)
		fmt.Printf("Room Number  : %d\n", b.RoomNumber)
		fmt.Printf("Check-In     : %s\n", b.CheckIn.Format("2006-01-02"))
		fmt.Printf("Check-Out    : %s\n", b.CheckOut.Format("2006-01-02"))
		fmt.Println("---------------")
	}
}

func (mh *ManagerHandler) UpdateEmployeeAvailability() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("\nEnter Employee Email: "))
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print(color.YellowString("Set Availability (true/false): "))
	var available bool
	fmt.Scanln(&available)

	err := mh.managerService.UpdateEmployeeAvailability(email, available)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	color.Green("Employee availability updated successfully.")
}

func (mh *ManagerHandler) ListEmployee() {
	color.Cyan("\n--- List of Employees ---")

	employees, err := mh.managerService.GetAllEmployees()
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	if len(employees) == 0 {
		color.Yellow("No employee found.")
		return
	}

	for _, emp := range employees {
		fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAvailable: %t\n---\n",
			emp.ID, emp.Name, emp.Email, emp.Role.String(), emp.Available)
	}
}

func (mh *ManagerHandler) DeleteEmployee() {
	color.Cyan("\n--- Delete Employee ---")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter Employee Email to delete: "))
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	err := mh.managerService.DeleteEmployeeByEmail(email)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	color.Green("Employee deleted successfully.")
}

func (h *ManagerHandler) AssignTasksToEmployees() {
	reader := bufio.NewReader(os.Stdin)

	requests, err := h.serviceRequestService.ViewUnassignedServiceRequest()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No unassigned service requests found.")
		return
	}

	color.Cyan("\n--- Unassigned Service Requests ---")
	for i, r := range requests {
		fmt.Printf("%d. Type: %s | Room: %d | Created At: %s\n", i+1, r.Type, r.RoomNum, r.CreatedAt)
	}

	fmt.Print(color.YellowString("Enter request numbers to assign (comma-separated): "))
	reqInput, _ := reader.ReadString('\n')
	reqInput = strings.TrimSpace(reqInput)
	reqParts := strings.Split(reqInput, ",")

	var selectedRequests []models.ServiceRequest
	for _, part := range reqParts {
		index, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || index < 1 || index > len(requests) {
			color.Red("Skipping invalid selection: %s", part)
			continue
		}
		selectedRequests = append(selectedRequests, requests[index-1])
	}

	if len(selectedRequests) == 0 {
		color.Yellow("No valid service requests selected.")
		return
	}

	for _, sr := range selectedRequests {
		color.Cyan("\n--- Assigning Request: %s (Room %d) ---", sr.Type, sr.RoomNum)

		bookingID, err := h.bookingService.GetBookingIDByRoomNumber(sr.RoomNum)
		if err != nil || bookingID == "" {
			color.Red("No booking found for room %d. Skipping...", sr.RoomNum)
			continue
		}

		staffList, err := h.managerService.GetAvailableStaffByTaskType(string(sr.Type))
		if err != nil || len(staffList) == 0 {
			color.Red("No available staff for this task. Skipping...")
			continue
		}

		color.Cyan("Available Staff:")
		if len(staffList) == 0 {
			fmt.Printf(color.RedString("No Employee Available for %s task."), sr.Type)
		}
		for i, s := range staffList {
			fmt.Printf("%d. %s (%s)\n", i+1, s.Name, s.Email)
		}

		fmt.Print(color.YellowString("Select staff by number: "))
		var staffChoice int
		fmt.Scanln(&staffChoice)
		if staffChoice < 1 || staffChoice > len(staffList) {
			color.Red("Invalid selection. Skipping...")
			continue
		}
		selectedStaff := staffList[staffChoice-1]

		fmt.Print(color.YellowString("Enter task details: "))
		details, _ := reader.ReadString('\n')
		details = strings.TrimSpace(details)

		err = h.managerService.AssignTaskFromServiceRequest(sr.ID, bookingID, details, selectedStaff.ID)
		if err != nil {
			color.Red("Failed to assign task: %v", err)
		} else {
			color.Green("Task assigned successfully.")
		}
	}
}

func (mh *ManagerHandler) ViewUnassignedServiceRequests() {
	color.Cyan("\n--- Unassigned Service Requests ---")

	requests, err := mh.serviceRequestService.ViewUnassignedServiceRequest()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No unassigned service requests found.")
		return
	}

	for _, req := range requests {
		fmt.Println("Request ID:", req.ID)
		fmt.Println("Guest ID:", req.UserID)
		fmt.Println("Room Number:", req.RoomNum)
		fmt.Println("Request Type:", req.Type)
		fmt.Println("Created At:", req.CreatedAt)
		fmt.Println("-----------------------------------")
	}
}

func (mh *ManagerHandler) ViewAllServiceRequests() {
	color.Cyan("\n--- All Guest Service Requests ---")

	requests, err := mh.serviceRequestService.ViewAllServiceRequests()
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No service requests found.")
		return
	}

	for _, req := range requests {
		fmt.Println("Request ID:", req.ID)
		fmt.Println("Guest ID:", req.UserID)
		fmt.Println("Room Number:", req.RoomNum)
		fmt.Println("Type:", req.Type)
		fmt.Println("Created At:", req.CreatedAt)
		fmt.Println("-----------")
	}
}

func (h *ManagerHandler) GenerateReport() {
	color.Cyan("\n--- Generating Hotel Report ---")

	err := h.managerService.PrintHotelReport()
	if err != nil {
		color.Red("Error generating report: %v", err)
	}
}

func (mh *ManagerHandler) ViewFeedback(ctx context.Context) {
	feedbacks, err := mh.managerService.ViewAllFeedback()
	if err != nil {
		color.Red("Error fetching feedback: %v", err)
		return
	}

	if len(feedbacks) == 0 {
		color.Yellow("No feedback available.")
		return
	}

	color.Cyan("\n--- All Feedback ---")
	if len(feedbacks) == 0 {
		fmt.Print(color.RedString("No Feedbacks."))
	}
	for _, fb := range feedbacks {
		fmt.Printf("Feedback ID: %s\nUser ID: %s\nMessage: %s\nDate: %s\n\n",
			fb.ID, fb.UserID, fb.Message, fb.CreatedAt)
	}
}

func (h *ManagerHandler) roomManagementMenu() {
RoomMgmtLoop:
	for {
		fmt.Println(titleStyle(config.RoomMgmtTitle))
		fmt.Println(optionStyle("\n1.") + " List Rooms")
		fmt.Println(optionStyle("2.") + " Add Room")
		fmt.Println(optionStyle("3.") + " Update Room")
		fmt.Println(optionStyle("4.") + " Delete Room")
		fmt.Println(optionStyle("5.") + " Back")
		fmt.Print(promptStyle(config.SelectOption))

		var rchoice int
		fmt.Scanln(&rchoice)
		switch rchoice {
		case 1:
			h.ListRooms()
		case 2:
			h.AddRoom()
		case 3:
			h.UpdateRoom()
		case 4:
			h.DeleteRoom()
		case 5:
			break RoomMgmtLoop
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (h *ManagerHandler) employeeManagementMenu() {
EmpMgmtLoop:
	for {
		fmt.Println(titleStyle(config.EmpMgmtTitle))
		fmt.Println(optionStyle("\n1.") + " List Employee")
		fmt.Println(optionStyle("2.") + " Update Employee Availability")
		fmt.Println(optionStyle("3.") + " Delete Employee")
		fmt.Println(optionStyle("4.") + " Back")
		fmt.Print(promptStyle(config.SelectOption))

		var echoice int
		fmt.Scanln(&echoice)
		switch echoice {
		case 1:
			h.ListEmployee()
		case 2:
			h.UpdateEmployeeAvailability()
		case 3:
			h.DeleteEmployee()
		case 4:
			break EmpMgmtLoop
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (h *ManagerHandler) serviceRequestManagementMenu() {
ServiceReqLoop:
	for {
		fmt.Println(titleStyle("Manage Service Requests"))
		fmt.Println(optionStyle("\n1.") + " View All Guest Service Requests")
		fmt.Println(optionStyle("2.") + " View Unassigned Guest Service Requests")
		fmt.Println(optionStyle("3.") + " Cancel Service Request")
		fmt.Println(optionStyle("4.") + " Back")
		fmt.Print(promptStyle(config.SelectOption))

		var sChoice int
		fmt.Scanln(&sChoice)
		switch sChoice {
		case 1:
			h.ViewAllServiceRequests()
		case 2:
			h.ViewUnassignedServiceRequests()
		case 3:
			h.CancelServiceRequest()
		case 4:
			break ServiceReqLoop
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (m *ManagerHandler) CancelServiceRequest() {
	color.Cyan("\n--- Unassigned Service Requests ---")

	requests, err := m.serviceRequestService.ViewUnassignedServiceRequest()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No unassigned service requests found.")
		return
	}

	for i, req := range requests {
		fmt.Printf("%d. Room Number: %d | Request Type: %s | Created At: %s\n",
			i+1, req.RoomNum, req.Type, req.CreatedAt)
	}

	var roomNum int
	fmt.Print(promptStyle("Enter Room Number for the service request to cancel: "))
	fmt.Scanln(&roomNum)

	err = m.serviceRequestService.CancelServiceRequestByRoomNum(roomNum)
	if err != nil {
		fmt.Println(errStyle("Error canceling service request:"), err)
		return
	}
	fmt.Println(successStyle("Service request canceled successfully."))
}

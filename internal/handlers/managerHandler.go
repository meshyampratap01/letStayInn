package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	managerservice "github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type ManagerHandler struct {
	roomService           roomService.IRoomService
	bookingService        bookingService.IBookingService
	userService           userService.IUserService
	serviceRequestService servicerequest.IServiceRequestService
	managerService        managerservice.IManagerService
}

func NewManagerHandler(rs roomService.IRoomService, bs bookingService.IBookingService, us userService.IUserService, srs servicerequest.IServiceRequestService, ms managerservice.IManagerService) *ManagerHandler {
	return &ManagerHandler{
		roomService:           rs,
		bookingService:        bs,
		userService:           us,
		serviceRequestService: srs,
		managerService:        ms,
	}
}

func (mh ManagerHandler) ManagerDashboardSummary() {
	totalRooms, err := mh.roomService.GetTotalRooms()
	if err != nil {
		color.Red("Error fetching total rooms: %v", err)
		return
	}

	availableRooms, err := mh.roomService.GetTotalAvailableRooms()
	if err != nil {
		color.Red("Error fetching available rooms: %v", err)
		return
	}

	bookedRooms, err := mh.bookingService.GetActiveBookings()
	if err != nil {
		color.Red("Error fetching booked rooms: %v", err)
		return
	}

	totalEmp, err := mh.managerService.GetTotalEmployees()
	if err != nil {
		color.Red("Error fetching total employees: %v", err)
		return
	}

	pendingRequests, err := mh.serviceRequestService.GetPendingRequestCount()
	if err != nil {
		color.Red("Error fetching pending service requests: %v", err)
		return
	}

	color.Cyan("\n--- Dashboard Summary ---")
	fmt.Printf("Total Rooms: %d\n", totalRooms)
	fmt.Printf("Booked Rooms: %d\n", len(bookedRooms))
	fmt.Printf("Available Rooms: %d\n", availableRooms)
	fmt.Printf("Total Employees: %d\n", totalEmp)
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

	fmt.Printf("%-10s %-15s %-12s %-12s %-30s\n",
		"Number", "Type", "Price (Rs)", "Available", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, room := range rooms {
		availability := color.RedString("No")
		if room.IsAvailable {
			availability = color.GreenString("Yes")
		}

		fmt.Printf("%-10d %-15s %-12.2f %-20s %-30s\n",
			room.Number,
			room.Type,
			room.Price,
			availability,
			room.Description)
	}
}

func (h *ManagerHandler) AddRoom() {
	reader := bufio.NewReader(os.Stdin)

	var number int
	for {
		fmt.Print(color.YellowString("Enter room number: "))
		_, err := fmt.Scanln(&number)
		if err != nil || number <= 0 {
			color.Red("Invalid room number. Please enter a positive number.")
			reader.ReadString('\n')
			continue
		}

		exists, err := h.roomService.RoomExists(number)
		if err != nil {
			color.Red("Error checking room number: %v", err)
			continue
		}
		if exists {
			color.Red("Room number already exists. Please enter a different number.")
			continue
		}

		break
	}

	var typeChoice int
	var roomType string
	for {
		fmt.Println(color.CyanString("Select room type:"))
		fmt.Println("1. Standard")
		fmt.Println("2. Deluxe")
		fmt.Println("3. Suite")
		fmt.Println("4. Executive")
		fmt.Print(color.YellowString("Enter choice (1-4): "))
		fmt.Scanln(&typeChoice)

		switch typeChoice {
		case 1:
			roomType = string(models.RoomTypeStandard)
		case 2:
			roomType = string(models.RoomTypeDeluxe)
		case 3:
			roomType = string(models.RoomTypeSuite)
		case 4:
			roomType = string(models.RoomTypeExecutive)
		default:
			color.Red("Invalid choice. Please try again.")
			continue
		}
		break
	}

	var price float64
	for {
		fmt.Print(color.YellowString("Enter room price: "))
		_, err := fmt.Scanln(&price)
		if err != nil || price <= 0 {
			color.Red("Invalid price. Please enter a positive number.")
			reader.ReadString('\n')
			continue
		}
		break
	}

	var availChoice int
	var isAvailable bool
	for {
		fmt.Println(color.CyanString("Is the room available?"))
		fmt.Println("1. Yes")
		fmt.Println("2. No")
		fmt.Print(color.YellowString("Enter choice: "))
		fmt.Scanln(&availChoice)

		switch availChoice {
		case 1:
			isAvailable = true
		case 2:
			isAvailable = false
		default:
			color.Red("Invalid choice. Please enter 1 or 2.")
			continue
		}
		break
	}

	fmt.Print(color.YellowString("Enter room description: "))
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

	booked, err := h.bookingService.IsRoomBooked(number)
	if err != nil {
		color.Red("Error checking booking status: %v", err)
		return
	}
	if booked {
		color.Red("Cannot update room %d — it is currently booked.", number)
		return
	}

	fmt.Print(color.CyanString("\nSelect what you want to update:\n"))
	fmt.Println("1. Room Type")
	fmt.Println("2. Price")
	fmt.Println("3. Availability")
	fmt.Println("4. Description")
	fmt.Print(color.YellowString("Enter choice: "))
	var choice int
	fmt.Scanln(&choice)

	var roomType, description string
	var price float64
	var isAvailable bool

	switch choice {
	case 1:
		fmt.Println(color.CyanString("Select new room type:"))
		fmt.Println("1. Standard")
		fmt.Println("2. Deluxe")
		fmt.Println("3. Suite")
		fmt.Println("4. Executive")
		fmt.Print(color.YellowString("Enter choice (1-4): "))
		var typeChoice int
		fmt.Scanln(&typeChoice)
		switch typeChoice {
		case 1:
			roomType = string(models.RoomTypeStandard)
		case 2:
			roomType = string(models.RoomTypeDeluxe)
		case 3:
			roomType = string(models.RoomTypeSuite)
		case 4:
			roomType = string(models.RoomTypeExecutive)
		default:
			color.Red("Invalid room type.")
			return
		}
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

	err = h.roomService.UpdateRoom(number, choice, roomType, price, isAvailable, description)
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

	booked, err := h.bookingService.IsRoomBooked(number)
	if err != nil {
		color.Red("Error checking booking status: %v", err)
		return
	}
	if booked {
		color.Red("Cannot delete room %d — it is currently booked.", number)
		return
	}

	err = h.roomService.DeleteRoom(number)
	if err != nil {
		color.Red("Error deleting room: %v", err)
		return
	}

	color.Green("Room deleted successfully.")
}

func (h *ManagerHandler) ListBookingsAndGuests() {
	bookings, err := h.bookingService.GetActiveBookings()
	if err != nil {
		color.Red("Error fetching bookings: %v", err)
		return
	}

	color.Cyan("\n--- Bookings and Guests ---")
	if len(bookings) == 0 {
		color.Yellow("No bookings found.")
		return
	}

	fmt.Printf("%-15s %-20s %-12s %-12s %-12s\n",
		"Room Number", "Guest Name", "Check-In", "Check-Out", "Status")
	fmt.Println(strings.Repeat("-", 85))

	for _, b := range bookings {
		guest, err := h.userService.GetUserNameByID(b.UserID)
		if err != nil {
			guest = "Unknown"
		}

		status := color.GreenString(b.Status)
		if strings.ToLower(b.Status) != "confirmed" {
			status = color.YellowString(b.Status)
		}

		fmt.Printf("%-15d %-20s %-12s %-12s %-12s \n",
			b.RoomNum,
			guest,
			b.CheckIn.Format("2006-01-02"),
			b.CheckOut.Format("2006-01-02"),
			status,
		)
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

func (h *ManagerHandler) GenerateReport() {
	color.Cyan("\n--- Generating Hotel Report ---")

	err := h.managerService.PrintHotelReport()
	if err != nil {
		color.Red("Error generating report: %v", err)
	}
}

func (mh *ManagerHandler) ViewFeedback() {
	feedbacks, err := mh.managerService.ViewAllFeedback()
	if err != nil {
		color.Red("❌ Error fetching feedback: %v", err)
		return
	}

	if len(feedbacks) == 0 {
		color.Yellow("⚠️ No feedback available.")
		return
	}

	color.Cyan("\n📋 --- All Feedback ---")
	for _, fb := range feedbacks {
		rating := "N/A"
		if fb.Rating > 0 {
			switch {
			case fb.Rating >= 4:
				rating = color.GreenString("%d/5", fb.Rating)
			case fb.Rating == 3:
				rating = color.YellowString("%d/5", fb.Rating)
			default:
				rating = color.RedString("%d/5", fb.Rating)
			}
		}

		room := "N/A"
		if fb.RoomNum != 0 {
			room = fmt.Sprintf("%d", fb.RoomNum)
		}

		booking := "N/A"
		if fb.BookingID != "" {
			booking = fb.BookingID
		}

		fmt.Printf(
			"\n🆔 Feedback ID : %s\n👤 User Name   : %s\n🔑 User ID     : %s\n💬 Message     : %s\n📅 Date        : %s\n🚪 Room Number : %s\n📖 Booking ID  : %s\n⭐ Rating      : %s\n",
			fb.ID,
			fb.UserName,
			fb.UserID,
			fb.Message,
			fb.CreatedAt.Format("2006-01-02 15:04"),
			room,
			booking,
			rating,
		)
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
		fmt.Println(optionStyle("2.") + " Create Employee")
		fmt.Println(optionStyle("3.") + " Delete Employee")
		fmt.Println(optionStyle("4.") + " Update Employee Availability")
		fmt.Println(optionStyle("5.") + " Back")
		fmt.Print(promptStyle(config.SelectOption))

		var echoice int
		fmt.Scanln(&echoice)
		switch echoice {
		case 1:
			h.ListEmployee()
		case 2:
			h.CreateEmployee()
		case 3:
			h.DeleteEmployee()
		case 4:
			h.UpdateEmployeeAvailability()
		case 5:
			break EmpMgmtLoop
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (h *ManagerHandler) CreateEmployee() {
	var email, password string
	var roleInt, availChoice int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Create Employee ---")
	fmt.Print("Enter Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	for {
		fmt.Print("Enter Email: ")
		fmt.Scanln(&email)

		if err := validators.ValidateEmail(email); err != nil {
			color.Red(err.Error())
			continue
		}
		break
	}

	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	fmt.Println("Select Role: ")
	fmt.Println("1. Kitchen Staff")
	fmt.Println("2. Cleaning Staff")
	fmt.Println("3. Manager")
	fmt.Print("Enter choice: ")
	fmt.Scanln(&roleInt)

	var role models.Role
	switch roleInt {
	case 1:
		role = models.RoleKitchenStaff
	case 2:
		role = models.RoleCleaningStaff
	case 3:
		role = models.RoleManager
	default:
		fmt.Println("Invalid role selection.")
		return
	}

	fmt.Println("Is Employee Available?")
	fmt.Println("1. Yes")
	fmt.Println("2. No")
	fmt.Print("Enter choice: ")
	fmt.Scanln(&availChoice)

	var available bool
	switch availChoice {
	case 1:
		available = true
	case 2:
		available = false
	default:
		fmt.Println("Invalid availability selection.")
		return
	}

	emp, err := h.userService.CreateEmployee(name, email, password, role, available)
	if err != nil {
		color.Red("Error creating employee: %v", err)
		return
	}

	color.Green("Employee created successfully! ID: %s", emp.ID)
}

func (h *ManagerHandler) serviceRequestManagementMenu() {
ServiceReqLoop:
	for {
		fmt.Println(titleStyle("\n--- Manage Service Requests ---"))
		fmt.Println(optionStyle("\n1.") + " View Pending / Unassigned Requests")
		fmt.Println(optionStyle("2.") + " Assign Request to Employee")
		fmt.Println(optionStyle("3.") + " Update Request Status")
		fmt.Println(optionStyle("4.") + " Cancel Service Request")
		fmt.Println(optionStyle("5.") + " Back")
		fmt.Print(promptStyle(config.SelectOption))

		var sChoice int
		fmt.Scanln(&sChoice)
		switch sChoice {
		case 1:
			h.ViewUnassignedServiceRequests()
		case 2:
			h.AssignServiceRequestToEmployee()
		case 3:
			h.UpdateServiceRequestStatus()
		case 4:
			h.CancelServiceRequest()
		case 5:
			break ServiceReqLoop
		default:
			fmt.Println(errStyle(config.InvalidOption))
		}
	}
}

func (mh *ManagerHandler) ViewUnassignedServiceRequests() {
	color.Cyan("\n--- Pending / Unassigned Service Requests ---")

	requests, err := mh.serviceRequestService.GetUnassignedServiceRequest()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No unassigned service requests found.")
		return
	}

	fmt.Printf("%-15s %-10s %-12s %-20s %-12s\n", "Guest Name", "Room No", "Type", "Created At", "Is Assigned")
	fmt.Println(strings.Repeat("-", 75))

	for _, req := range requests {
		guestName, err := mh.userService.GetUserNameByID(req.UserID)
		if err != nil {
			color.Red("Error fetching guest name for UserID %s: %v", req.UserID, err)
			continue
		}

		isAssignedStr := color.RedString("No")
		if req.IsAssigned {
			isAssignedStr = color.GreenString("Yes")
		}

		fmt.Printf("%-15s %-10d %-12s %-20s %-12s\n",
			guestName,
			req.RoomNum,
			req.Type,
			req.CreatedAt.Format("02-01-2006 15:04"),
			isAssignedStr)
	}
}

func (m *ManagerHandler) CancelServiceRequest() {
	color.Cyan("\n--- Unassigned Service Requests ---")

	requests, err := m.serviceRequestService.GetUnassignedServiceRequest()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	if len(requests) == 0 {
		color.Yellow("No unassigned service requests found.")
		return
	}

	fmt.Printf("%-5s %-10s %-15s %-20s\n", "No.", "Room No", "Type", "Created At")
	fmt.Println(strings.Repeat("-", 60))
	for i, req := range requests {
		fmt.Printf("%-5d %-10d %-15s %-20s\n",
			i+1, req.RoomNum, req.Type, req.CreatedAt.Format("2006-01-02 15:04"))
	}

	var choice int
	fmt.Print(promptStyle("Enter the number of the service request to cancel: "))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(requests) {
		fmt.Println(errStyle("Invalid selection"))
		return
	}

	selectedReq := requests[choice-1]

	err = m.serviceRequestService.CancelServiceRequestByID(selectedReq.ID)
	if err != nil {
		fmt.Println(errStyle("Error canceling service request:"), err)
		return
	}

	fmt.Println(successStyle("Service request canceled successfully."))
}

func (h *ManagerHandler) AssignServiceRequestToEmployee() {
	fmt.Println(titleStyle("Unassigned Service Requests\n"))
	unassigned, err := h.serviceRequestService.GetUnassignedServiceRequest()
	if err != nil {
		fmt.Println(errStyle("Error fetching unassigned requests:", err))
		return
	}

	if len(unassigned) == 0 {
		color.Yellow("No pending / unassigned service requests found.")
		return
	}

	fmt.Printf("%-5s %-10s %-15s %-20s %-30s\n", "No.", "Room No", "Type", "Created At", "Details")
	fmt.Println(strings.Repeat("-", 85))
	for i, r := range unassigned {
		fmt.Printf("%-5d %-10d %-15s %-20s %-30s\n",
			i+1, r.RoomNum, r.Type, r.CreatedAt.Format("2006-01-02 15:04"), r.Details)
	}

	var choice int
	fmt.Print(promptStyle("Enter the number of the request to assign: "))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(unassigned) {
		fmt.Println(errStyle("Invalid selection"))
		return
	}

	selectedReq := unassigned[choice-1]

	var requiredRole models.Role
	switch selectedReq.Type {
	case models.ServiceTypeCleaning:
		requiredRole = models.RoleCleaningStaff
	case models.ServiceTypeFood:
		requiredRole = models.RoleKitchenStaff
	default:
		fmt.Println(errStyle("Unknown service request type. Cannot assign."))
		return
	}

	employees, err := h.userService.GetAllEmployees()
	if err != nil {
		fmt.Println(errStyle("Error fetching employees:", err))
		return
	}
	var filteredEmps []models.User
	for _, u := range employees {
		if u.Role == requiredRole {
			filteredEmps = append(filteredEmps, u)
		}
	}
	if len(filteredEmps) == 0 {
		color.Yellow("No employees for this request type.")
		return
	}

	for {
		fmt.Printf("%-5s %-20s %-20s %-12s\n", "No.", "Name", "Email", "Available")
		for i, emp := range filteredEmps {
			availStr := color.GreenString("Yes")
			if !emp.Available {
				availStr = color.RedString("No")
			}
			fmt.Printf("%-5d %-20s %-20s %-12s\n", i+1, emp.Name, emp.Email, availStr)
		}
		var empChoice int
		fmt.Print(promptStyle("Select employee number to assign: "))
		fmt.Scanln(&empChoice)
		if empChoice < 1 || empChoice > len(filteredEmps) {
			fmt.Println(errStyle("Invalid employee selection"))
			continue
		}
		emp := filteredEmps[empChoice-1]
		if !emp.Available {
			color.Red("Selected employee is not available. Please choose an available employee.")
			continue
		}
		err = h.managerService.AssignServiceRequest(selectedReq.ID, emp.ID)
		if err != nil {
			fmt.Println(errStyle("Error assigning request:", err))
			return
		}
		err = h.serviceRequestService.UpdateServiceRequestAssignment(selectedReq.ID, true)
		if err != nil {
			fmt.Println(errStyle("Error marking request as assigned:", err))
			return
		}
		fmt.Println(successStyle("Service request assigned successfully!"))
		break
	}
}

func (h *ManagerHandler) UpdateServiceRequestStatus() {
	fmt.Println(titleStyle("Update Service Request Status"))

	requests, err := h.serviceRequestService.GetUnassignedServiceRequest()
	if err != nil {
		fmt.Println(errStyle(fmt.Sprintf("Error fetching service requests: %v", err)))
		return
	}

	if len(requests) == 0 {
		fmt.Println(errStyle("No service requests found."))
		return
	}

	fmt.Printf("%-10s %-15s %-20s %-15s %-30s\n", "Room No", "Type", "Status", "Assigned", "Details")
	fmt.Println(strings.Repeat("-", 90))
	for _, req := range requests {
		fmt.Printf("%-10d %-15s %-20s %-15t %-30s\n",
			req.RoomNum,
			string(req.Type),
			string(req.Status),
			req.IsAssigned,
			req.Details,
		)
	}

	var roomNum int
	fmt.Print(promptStyle("\nEnter Room Number to update: "))
	fmt.Scanln(&roomNum)

	var targetRequest *models.ServiceRequest
	for _, req := range requests {
		if req.RoomNum == roomNum {
			targetRequest = &req
			break
		}
	}

	if targetRequest == nil {
		fmt.Println(errStyle("No service request found for that room number."))
		return
	}

	fmt.Println(optionStyle("\n1.") + " Pending")
	fmt.Println(optionStyle("2.") + " In Progress")
	fmt.Println(optionStyle("3.") + " Done")
	fmt.Println(optionStyle("4.") + " Cancelled")
	fmt.Print(promptStyle("Select new status: "))

	var choice int
	fmt.Scanln(&choice)

	var newStatus models.ServiceStatus
	switch choice {
	case 1:
		newStatus = models.ServiceStatusPending
	case 2:
		newStatus = models.ServiceStatusInProgress
	case 3:
		newStatus = models.ServiceStatusDone
	case 4:
		newStatus = models.ServiceStatusCancelled
	default:
		fmt.Println(errStyle("Invalid status option."))
		return
	}

	err = h.serviceRequestService.UpdateServiceRequestStatus(targetRequest.ID, newStatus)
	if err != nil {
		fmt.Println(errStyle(fmt.Sprintf("Error updating status: %v", err)))
		return
	}

	fmt.Println(successStyle("Service request status updated successfully!"))
}

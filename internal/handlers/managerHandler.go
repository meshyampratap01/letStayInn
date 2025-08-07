package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	fmt.Println("\n--- Dashboard Summary ---")
	fmt.Printf("Total Rooms: %d\n", totalRooms)
	fmt.Printf("Booked Rooms: %d\n", bookedRooms)
	fmt.Printf("Available Rooms: %d\n", availableRooms)
	fmt.Printf("Total Guests: %d\n", totalGuests)
	fmt.Printf("Total Employee: %d\n", totalEmp)
	fmt.Printf("Pending Service Requests: %d\n", pendingRequests)
}

func (h *ManagerHandler) ListRooms() {
	rooms, err := h.roomService.GetAllRooms()
	if err != nil {
		fmt.Println("Error fetching rooms:", err)
		return
	}

	if len(rooms) == 0 {
		fmt.Println("No rooms found.")
		return
	}

	fmt.Println("\n--- Room List ---")
	for _, room := range rooms {
		availability := "No"
		if room.IsAvailable {
			availability = "Yes"
		}
		fmt.Printf("Number: %d | Type: %s | Price: %.2f | Available: %s | Description: %s\n",
			room.Number, room.Type, room.Price, availability, room.Description)
	}
}

func (h *ManagerHandler) AddRoom() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter room number: ")
	var number int
	fmt.Scanln(&number)

	fmt.Print("Enter room type (e.g., Deluxe, Standard): ")
	roomType, _ := reader.ReadString('\n')
	roomType = strings.TrimSpace(roomType)

	fmt.Print("Enter room price: ")
	var price float64
	fmt.Scanln(&price)

	fmt.Print("Is the room available? (yes/no): ")
	availableInput, _ := reader.ReadString('\n')
	availableInput = strings.TrimSpace(strings.ToLower(availableInput))
	isAvailable := availableInput == "yes"

	fmt.Print("Enter room description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	err := h.roomService.AddRoom(number, roomType, price, isAvailable, description)
	if err != nil {
		fmt.Println("Error adding room:", err)
		return
	}

	fmt.Println("Room added successfully.")
}

func (h *ManagerHandler) UpdateRoom() {
	reader := bufio.NewReader(os.Stdin)

	h.ListRooms()

	fmt.Print("\nEnter Room Number to update: ")
	var number int
	fmt.Scanln(&number)

	fmt.Println("\nSelect what you want to update:")
	fmt.Println("1. Room Type")
	fmt.Println("2. Price")
	fmt.Println("3. Availability")
	fmt.Println("4. Description")
	fmt.Print("Enter choice: ")
	var choice int
	fmt.Scanln(&choice)

	var roomType, description string
	var price float64
	var isAvailable bool

	switch choice {
	case 1:
		fmt.Print("Enter new Room Type: ")
		roomType, _ = reader.ReadString('\n')
		roomType = strings.TrimSpace(roomType)
	case 2:
		fmt.Print("Enter new Price: ")
		fmt.Scanln(&price)
	case 3:
		fmt.Print("Is room available? (yes/no): ")
		availableInput, _ := reader.ReadString('\n')
		isAvailable = strings.TrimSpace(strings.ToLower(availableInput)) == "yes"
	case 4:
		fmt.Print("Enter new Description: ")
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
	default:
		fmt.Println("Invalid choice.")
		return
	}

	err := h.roomService.UpdateRoom(number, choice, roomType, price, isAvailable, description)
	if err != nil {
		fmt.Println("Error updating room:", err)
		return
	}

	fmt.Println("Room updated successfully.")
}

func (h *ManagerHandler) DeleteRoom() {
	h.ListRooms()
	fmt.Print("\nEnter Room Number to delete: ")
	var number int
	fmt.Scanln(&number)

	err := h.roomService.DeleteRoom(number)
	if err != nil {
		fmt.Println("Error deleting room:", err)
		return
	}

	fmt.Println("Room deleted successfully.")
}

func (h *ManagerHandler) ListBookingsAndGuests() {
	bookings, err := h.bookingService.GetAllBookingsWithGuests()
	if err != nil {
		fmt.Println("Error fetching bookings:", err)
		return
	}

	fmt.Println("\n--- Bookings and Guests ---")
	if len(bookings) == 0 {
		fmt.Println("No bookings found.")
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

	fmt.Print("\nEnter Employee Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Set Availability (true/false): ")
	var available bool
	fmt.Scanln(&available)

	err := mh.managerService.UpdateEmployeeAvailability(email, available)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Employee availability updated successfully.")
}

func (mh *ManagerHandler) ListEmployee() {
	fmt.Println("\n--- List of Employees ---")

	employees, err := mh.managerService.GetAllEmployees()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(employees) == 0 {
		fmt.Println("No emplyee found.")
		return
	}

	for _, emp := range employees {
		fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAvailable: %t\n---\n",
			emp.ID, emp.Name, emp.Email, emp.Role.String(), emp.Available)
	}
}

func (mh *ManagerHandler) DeleteEmployee() {
	fmt.Println("\n--- Delete Employee ---")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Employee Email to delete: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	err := mh.managerService.DeleteEmployeeByEmail(email)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Employee deleted successfully.")
}

func (h *ManagerHandler) AssignTasksToEmployees() {
	reader := bufio.NewReader(os.Stdin)

	requests, err := h.serviceRequestService.ViewUnassignedServiceRequest()
	if err != nil {
		fmt.Println("Error fetching service requests:", err)
		return
	}

	if len(requests) == 0 {
		fmt.Println("No unassigned service requests found.")
		return
	}

	fmt.Println("\n--- Unassigned Service Requests ---")
	for i, r := range requests {
		fmt.Printf("%d. Type: %s | Room: %d | Created At: %s\n", i+1, r.Type, r.RoomNum, r.CreatedAt)
	}

	
	fmt.Print("Enter request numbers to assign (comma-separated): ")
	reqInput, _ := reader.ReadString('\n')
	reqInput = strings.TrimSpace(reqInput)
	reqParts := strings.Split(reqInput, ",")

	var selectedRequests []models.ServiceRequest
	for _, part := range reqParts {
		index, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || index < 1 || index > len(requests) {
			fmt.Printf("Skipping invalid selection: %s\n", part)
			continue
		}
		selectedRequests = append(selectedRequests, requests[index-1])
	}

	if len(selectedRequests) == 0 {
		fmt.Println("No valid service requests selected.")
		return
	}


	for _, sr := range selectedRequests {
		fmt.Printf("\n--- Assigning Request: %s (Room %d) ---\n", sr.Type, sr.RoomNum)


		bookingID, err := h.bookingService.GetBookingIDByRoomNumber(sr.RoomNum)
		if err != nil || bookingID == "" {
			fmt.Printf("No booking found for room %d. Skipping...\n", sr.RoomNum)
			continue
		}


		staffList, err := h.managerService.GetAvailableStaffByTaskType(string(sr.Type))
		if err != nil || len(staffList) == 0 {
			fmt.Println("No available staff for this task. Skipping...")
			continue
		}

		fmt.Println("Available Staff:")
		for i, s := range staffList {
			fmt.Printf("%d. %s (%s)\n", i+1, s.Name, s.Email)
		}

		fmt.Print("Select staff by number: ")
		var staffChoice int
		fmt.Scanln(&staffChoice)
		if staffChoice < 1 || staffChoice > len(staffList) {
			fmt.Println("Invalid selection. Skipping...")
			continue
		}
		selectedStaff := staffList[staffChoice-1]

		fmt.Print("Enter task details: ")
		details, _ := reader.ReadString('\n')
		details = strings.TrimSpace(details)


		err = h.managerService.AssignTaskFromServiceRequest(sr.ID, bookingID, details, selectedStaff.ID)
		if err != nil {
			fmt.Println("Failed to assign task:", err)
		} else {
			fmt.Println("Task assigned successfully.")
		}
	}
}

func (mh *ManagerHandler) ViewUnassignedServiceRequests() {
	fmt.Println("\n--- Unassigned Service Requests ---")

	requests, err := mh.serviceRequestService.ViewUnassignedServiceRequest()
	if err != nil {
		fmt.Println("Error fetching service requests:", err)
		return
	}

	if len(requests) == 0 {
		fmt.Println("No unassigned service requests found.")
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
	fmt.Println("\n--- All Guest Service Requests ---")

	requests, err := mh.serviceRequestService.ViewAllServiceRequests()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(requests) == 0 {
		fmt.Println("No service requests found.")
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
	fmt.Println("\n--- Generating Hotel Report ---")

	err := h.managerService.PrintHotelReport()
	if err != nil {
		fmt.Println("Error generating report:", err)
	}
}

package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
)

type EmployeeHandler struct {
	employeeService employeeService.IEmployeeService
}

func NewEmployeeHandler(employeeService employeeService.IEmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

func (eh *EmployeeHandler) ViewAssignedServiceRequests(employeeID string) error {
	requests, err := eh.employeeService.GetAssignedServiceRequests(employeeID)
	if err != nil {
		return fmt.Errorf("error fetching assigned requests: %v", err)
	}

	if len(requests) == 0 {
		color.Yellow("No service requests assigned.")
		return nil
	}

	color.Cyan("\n--- Assigned Service Requests ---\n")
	fmt.Printf("%-5s %-15s %-12s %-15s %-30s\n",
		"No", "Type", "Room No", "Status", "Details")
	fmt.Println(strings.Repeat("-", 80))

	for i, req := range requests {
		roomNum, err := eh.employeeService.GetRoomNumberByBookingID(req.BookingID)
		if err != nil {
			roomNum = color.RedString("Unknown")
		}

		fmt.Printf("%-5d %-15s %-12s %-15s %-30s\n",
			i+1,
			req.Type,
			roomNum,
			req.Status,
			req.Details,
		)
	}
	return nil
}


func (eh *EmployeeHandler) UpdateServiceRequestStatus(employeeID string) error {
	requests, err := eh.employeeService.GetAssignedServiceRequests(employeeID)
	if err != nil {
		return fmt.Errorf("error fetching service requests: %w", err)
	}

	if len(requests) == 0 {
		color.Yellow("No service requests assigned.")
		return nil
	}

	color.Cyan("\n--- Your Assigned Service Requests ---\n")
	fmt.Printf("%-5s %-15s %-12s %-15s %-30s\n",
		"No", "Type", "Room No", "Status", "Details")
	fmt.Println(strings.Repeat("-", 80))

	for i, req := range requests {
		roomNum, err := eh.employeeService.GetRoomNumberByBookingID(req.BookingID)
		if err != nil {
			roomNum = color.RedString("Unknown")
		}

		fmt.Printf("%-5d %-15s %-12s %-15s %-30s\n",
			i+1,
			req.Type,
			roomNum,
			req.Status,
			req.Details,
		)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("\nEnter the number of the request you want to update: "))
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(requests) {
		return fmt.Errorf("invalid request selection")
	}

	selectedRequest := requests[choice-1]

	fmt.Println(color.YellowString("\nSelect New Status:"))
	fmt.Println("1. Pending")
	fmt.Println("2. InProgress")
	fmt.Println("3. Done")

	fmt.Print(color.YellowString("Enter choice (1-3): "))
	statusChoice, _ := reader.ReadString('\n')
	statusChoice = strings.TrimSpace(statusChoice)

	var newStatus models.ServiceStatus
	switch statusChoice {
	case "1":
		newStatus = models.ServiceStatusPending
	case "2":
		newStatus = models.ServiceStatusInProgress
	case "3":
		newStatus = models.ServiceStatusDone
	default:
		return fmt.Errorf("invalid status choice")
	}

	if err := eh.employeeService.UpdateServiceRequestStatus(selectedRequest.ID, newStatus); err != nil {
		return fmt.Errorf("error updating request status: %w", err)
	}

	color.Green("Service request status updated successfully.")
	return nil
}


func (eh *EmployeeHandler) ToggleAvailability(userID string) error {
	available, err := eh.employeeService.GetAvailability(userID)
	if err != nil {
		return fmt.Errorf("error retrieving availability: %w", err)
	}

	status := color.RedString("Unavailable")
	if available {
		status = color.GreenString("Available")
	}
	color.Cyan("\n--- Toggle Availability ---")
	fmt.Println("Current status:", status)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Do you want to toggle your availability? (y/n): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" {
		color.Yellow("No changes made.")
		return nil
	}

	if err := eh.employeeService.ToggleAvailability(userID); err != nil {
		return fmt.Errorf("error toggling availability: %w", err)
	}

	newStatus := color.RedString("Unavailable")
	if !available {
		newStatus = color.GreenString("Available")
	}
	color.Green("Availability toggled successfully. New status: %s", newStatus)

	return nil
}


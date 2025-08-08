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

func (eh *EmployeeHandler) ViewAssignedTasks(userID string) {
	tasks, err := eh.employeeService.ViewAssignedTasks(userID)
	if err != nil {
		color.Red("Error fetching assigned tasks: %v", err)
		return
	}

	color.Cyan("\n--- Assigned Tasks ---")
	for i, task := range tasks {
		roomNum, err := eh.employeeService.GetRoomNumberByBookingID(task.BookingID)
		if err != nil {
			color.Red("Error: %v", err)
		}
		fmt.Printf("%d. Type: %s | Room Number: %s | Status: %s\n",
			i+1, task.Type, roomNum, task.Status)
	}
}

func (eh *EmployeeHandler) UpdateTaskStatus(userID string) {
	tasks, err := eh.employeeService.ViewAssignedTasks(userID)
	if err != nil {
		color.Red("Error fetching tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		color.Yellow("No tasks assigned.")
		return
	}

	color.Cyan("\n--- Your Assigned Tasks ---")
	for i, t := range tasks {
		roomNumber, err := eh.employeeService.GetRoomNumberByBookingID(t.BookingID)
		if err != nil {
			roomNumber = color.RedString("Unknown")
		}

		fmt.Printf("%d. Room: %s | Type: %s | Status: %s\n", i+1, roomNumber, t.Type, t.Status)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("\nEnter the number of the task you want to update: "))
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(tasks) {
		color.Red("Invalid task selection.")
		return
	}

	selectedTask := tasks[choice-1]

	fmt.Print(color.YellowString("Select New Status:"))
	fmt.Println("1. Pending")
	fmt.Println("2. InProgress")
	fmt.Println("3. Done")

	fmt.Print(color.YellowString("Enter choice (1-3): "))
	statusChoice, _ := reader.ReadString('\n')
	statusChoice = strings.TrimSpace(statusChoice)

	var newStatus models.TaskStatus
	switch statusChoice {
	case "1":
		newStatus = models.TaskStatusPending
	case "2":
		newStatus = models.TaskStatusInProgress
	case "3":
		newStatus = models.TaskStatusDone
	default:
		color.Red("Invalid status choice.")
		return
	}

	err = eh.employeeService.UpdateTaskStatus(selectedTask.ID, newStatus)
	if err != nil {
		color.Red("Error updating task status: %v", err)
		return
	}

	color.Green("Task status updated successfully.")
}

func (eh *EmployeeHandler) ToggleAvailability(userID string) {
	available, err := eh.employeeService.GetAvailability(userID)
	if err != nil {
		color.Red("Error retrieving availability: %v", err)
		return
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
		return
	}

	err = eh.employeeService.ToggleAvailability(userID)
	if err != nil {
		color.Red("Error toggling availability: %v", err)
		return
	}

	newStatus := color.RedString("Unavailable")
	if !available {
		newStatus = color.GreenString("Available")
	}
	color.Green("Availability toggled successfully. New status: %s", newStatus)
}

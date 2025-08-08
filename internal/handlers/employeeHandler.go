package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		fmt.Println("Error fetching assigned tasks:", err)
		return
	}

	fmt.Println("\n--- Assigned Tasks ---")
	for i, task := range tasks {
		roomNum,err:=eh.employeeService.GetRoomNumberByBookingID(task.BookingID)
		if err!=nil{
			fmt.Println("Error:",err)
		}
		fmt.Printf("%d. Type: %s | Room Number: %s| Status: %s\n",
			i+1,task.Type,roomNum, task.Status )
	}
}

func (eh *EmployeeHandler) UpdateTaskStatus(userID string) {
	tasks, err := eh.employeeService.ViewAssignedTasks(userID)
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks assigned.")
		return
	}

	fmt.Println("\n--- Your Assigned Tasks ---")
	for i, t := range tasks {
		roomNumber, err := eh.employeeService.GetRoomNumberByBookingID(t.BookingID)
		if err != nil {
			roomNumber = "Unknown"
		}

		fmt.Printf("%d. Room: %s | Type: %s | Status: %s\n", i+1, roomNumber, t.Type, t.Status)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter the number of the task you want to update: ")
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(tasks) {
		fmt.Println("Invalid task selection.")
		return
	}

	selectedTask := tasks[choice-1]

	fmt.Println("Select New Status:")
	fmt.Println("1. Pending")
	fmt.Println("2. InProgress")
	fmt.Println("3. Done")

	fmt.Print("Enter choice (1-3): ")
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
		fmt.Println("Invalid status choice.")
		return
	}

	err = eh.employeeService.UpdateTaskStatus(selectedTask.ID, newStatus)
	if err != nil {
		fmt.Println("Error updating task status:", err)
		return
	}

	fmt.Println("Task status updated successfully.")
}

func (eh *EmployeeHandler) ToggleAvailability(userID string) {

	available, err := eh.employeeService.GetAvailability(userID)
	if err != nil {
		fmt.Println("Error retrieving availability:", err)
		return
	}

	status := "Unavailable"
	if available {
		status = "Available"
	}
	fmt.Println("\n--- Toggle Availability ---")
	fmt.Println("Current status:", status)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to toggle your availability? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" {
		fmt.Println("No changes made.")
		return
	}

	err = eh.employeeService.ToggleAvailability(userID)
	if err != nil {
		fmt.Println("Error toggling availability:", err)
		return
	}

	newStatus := "Unavailable"
	if !available {
		newStatus = "Available"
	}
	fmt.Println("Availability toggled successfully. New status:", newStatus)
}

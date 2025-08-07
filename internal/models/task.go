package models

import "time"

type TaskType string
type TaskStatus string

const (
	// Task Types
	TaskTypeCleaning TaskType = "Cleaning"
	TaskTypeFood     TaskType = "Food"

	// Task Statuses
	TaskStatusPending    TaskStatus = "Pending"
	TaskStatusInProgress TaskStatus = "In Progress"
	TaskStatusDone       TaskStatus = "Done"
)

type Task struct {
	ID         string      `json:"id"`
	Type       TaskType    `json:"type"` // Cleaning, Food
	AssignedTo string      `json:"assigned_to"`
	BookingID  string      `json:"booking_id"`
	Status     TaskStatus  `json:"status"` // Pending, In Progress, Done
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Details    string      `json:"details"`
}

package models

import "time"

type ServiceType string
type ServiceStatus string

const (
	// Service Types
	ServiceTypeCleaning ServiceType = "Cleaning"
	ServiceTypeFood     ServiceType = "Food"

	// Service Statuses
	ServiceStatusPending    ServiceStatus = "Pending"
	ServiceStatusInProgress ServiceStatus = "In Progress"
	ServiceStatusDone       ServiceStatus = "Done"
	ServiceStatusCancelled  ServiceStatus = "Cancelled"
)

type ServiceRequest struct {
	ID          string        `json:"id"`            // Primary Key
	UserID      string        `json:"user_id"`       // FK → Users
	BookingID   string        `json:"booking_id"`    // FK → Bookings
	RoomNum     int           `json:"room_num"`      // Redundant but useful for quick lookup
	Type        ServiceType   `json:"type"`          // Cleaning / Food
	Status      ServiceStatus `json:"status"`        // Pending / In Progress / Done / Cancelled
	IsAssigned  bool          `json:"is_assigned"`   // Assigned or not
	AssignedTo  string        `json:"assigned_to"`   // FK → Employees (UserID of staff)
	Details     string        `json:"details"`       // Additional description
	CreatedAt   time.Time     `json:"created_at"`    // When request was created
	UpdatedAt   time.Time     `json:"updated_at"`    // Last update time
}

package models

type ServiceType string
type ServiceStatus string

const (
	// Service Types
	ServiceTypeCleaning ServiceType = "Cleaning"
	ServiceTypeFood     ServiceType = "Food"

	// Service Statuses
	ServiceStatusPending   ServiceStatus = "Pending"
	ServiceStatusInProgrss ServiceStatus = "In Progress"
	ServiceStatusDone      ServiceStatus = "Done"
	ServiceStatusCancelled ServiceStatus = "Cancelled"
)

type ServiceRequest struct {
	ID         string        `json:"id"`
	UserID     string        `json:"user_id"`
	BookingID  string        `json:"booking_id"`
	RoomNum    int           `json:"room_num"`
	Type       ServiceType   `json:"type"`
	Status     ServiceStatus `json:"status"`
	CreatedAt  string        `json:"created_at"`
	IsAssigned bool          `json:"is_assigned"`
	Details    string     		`json:"details"` 
}

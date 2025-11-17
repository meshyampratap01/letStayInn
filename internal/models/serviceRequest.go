package models

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
	PK         string        `dynamodbav:"pk" json:"pk"`
	SK         string        `dynamodbav:"sk" json:"sk"`
	ID         string        `dynamodbav:"id" json:"id"`
	UserID     string        `dynamodbav:"userID" json:"user_id"`
	BookingID  string        `dynamodbav:"bookingID" json:"booking_id"`
	RoomNum    int           `dynamodbav:"roomNum" json:"room_num"`
	Type       ServiceType   `dynamodbav:"type" json:"type"`
	Status     ServiceStatus `dynamodbav:"status" json:"status"`
	IsAssigned bool          `dynamodbav:"isAssigned" json:"is_assigned"`
	AssignedTo string        `dynamodbav:"assignedTo" json:"assigned_to"`
	Details    string        `dynamodbav:"details" json:"details"`
}

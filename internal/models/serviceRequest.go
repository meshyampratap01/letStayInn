package models

type ServiceType string
type ServiceStatus string

const (
	// Service Types
	ServiceTypeCleaning ServiceType = "Cleaning"
	ServiceTypeFood     ServiceType = "Food"

	// Service Statuses
	ServiceStatusPending  ServiceStatus = "Pending"
	ServiceStatusAccepted ServiceStatus = "Accepted"
	ServiceStatusDone     ServiceStatus = "Done"
)

type ServiceRequest struct {
	ID        	string   				`json:"id"`
	UserID    	string   				`json:"user_id"`
	RoomNum		int						`json:"room_num"`
	Type      	ServiceType				`json:"type"` 
	Status    	ServiceStatus  			`json:"status"`
	CreatedAt 	string  				`json:"created_at"`
	IsAssigned  bool      				`json:"is_assigned"`
}
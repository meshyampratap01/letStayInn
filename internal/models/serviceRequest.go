package models

import (
	"time"

	"gorm.io/gorm"
)

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
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"` // Primary Key
	UserID     string         `gorm:"type:uuid;not null" json:"user_id"`                        // FK → Users
	BookingID  string         `gorm:"type:uuid" json:"booking_id"`                              // FK → Bookings
	RoomNum    int            `json:"room_num"`                                                 // Redundant but useful for quick lookup
	Type       ServiceType    `json:"type"`                                                     // Cleaning / Food
	Status     ServiceStatus  `json:"status"`                                                   // Pending / In Progress / Done / Cancelled
	IsAssigned bool           `json:"is_assigned"`                                              // Assigned or not
	AssignedTo string         `gorm:"type:uuid;default:null;constraint:OnDelete:SET NULL;" json:"assigned_to"`                             // FK → Employees (UserID of staff)
	Details    string         `json:"details"`                                                  // Additional description
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID"`
	Booking Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:SET NULL;"`
	Staff   User    `gorm:"foreignKey:AssignedTo;constraint:OnDelete:SET NULL;"`
}

package models

import (
	"fmt"
	"strings"
)

type Role int

const (
	RoleGuest Role = iota + 1
	RoleKitchenStaff
	RoleCleaningStaff
	RoleManager
)

func (r Role) String() string {
	switch r {
	case RoleGuest:
		return "Guest"
	case RoleKitchenStaff:
		return "KitchenStaff"
	case RoleCleaningStaff:
		return "CleaningStaff"
	case RoleManager:
		return "Manager"
	default:
		return "Unknown"
	}
}

func (r Role) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *Role) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	switch str {
	case "Guest":
		*r = RoleGuest
	case "KitchenStaff":
		*r = RoleKitchenStaff
	case "CleaningStaff":
		*r = RoleCleaningStaff
	case "Manager":
		*r = RoleManager
	default:
		return fmt.Errorf("invlid role: %s", str)
	}
	return nil
}

type User struct {
	PK        string    `dynamodbav:"pk" json:"pk"`
	SK        string    `dynamodbav:"sk" json:"sk"`
	ID        string    `dynamodbav:"id" json:"id"`
	Name      string    `dynamodbav:"name" json:"name"`
	Email     string    `dynamodbav:"email" json:"email"`
	Password  string    `dynamodbav:"password" json:"password"`
	Role      Role      `dynamodbav:"role" json:"role"`
	Available bool      `dynamodbav:"available" json:"available"` 

	// Relationships
	Bookings        []Booking        `json:"-"`
	Feedbacks       []Feedback       `json:"-"`
	ServiceRequests []ServiceRequest `json:"-"`
}

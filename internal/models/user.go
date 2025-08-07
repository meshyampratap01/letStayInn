package models

import (
	"fmt"
	"strings"
	"time"
)

type Role int

const(
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

func (r Role) MarshalJSON() ([]byte,error){
	return []byte(`"`+r.String()+`"`),nil
}

func (r *Role) UnmarshalJSON(data []byte) error{
	str:=strings.Trim(string(data), `"`)
	switch str{
	case "Guest":
		*r=RoleGuest
	case "KitchenStaff":
		*r=RoleKitchenStaff
	case "CleaningStaff":
		*r=RoleCleaningStaff
	case "Manager":
		*r= RoleManager
	default:
		return fmt.Errorf("invlid role: %s",str)
	}
	return nil
} 


type User struct{
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      Role   	`json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Available bool      `json:"available"` // for staff
}
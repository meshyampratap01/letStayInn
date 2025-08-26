package dto

type ServiceRequestDTO struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	RoomNum    int `json:"room_num"`
	Type       string `json:"type"`
	Details    string `json:"details"`
	CreatedAt  string `json:"created_at"`
	IsAssigned bool   `json:"is_assigned"`
	Status     string `json:"status"`
	EmployeeID string `json:"employee_id,omitempty"`
}

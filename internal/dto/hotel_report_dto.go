package dto

type HotelReportDTO struct {
	TotalRooms            int            `json:"total_rooms"`
	AvailableRooms        int            `json:"available_rooms"`
	TotalStaff            int            `json:"total_staff"`
	TotalBookings         int            `json:"total_bookings"`
	UnassignedRequests    int            `json:"unassigned_requests"`
	ServiceRequestSummary map[string]int `json:"service_request_summary"`
}

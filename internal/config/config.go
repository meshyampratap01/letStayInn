package config

const (
	UsersFile    = "data/users.json"
	RoomsFile    = "data/rooms.json"
	BookingsFile = "data/bookings.json"
	TasksFile    = "data/tasks.json"
	ServiceRequestFile = "data/serviceRequests.json"
	FeedbackFile = "data/feedbacks.json"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Panic("Error loading .env file")
// 	}

// 	RoomsFile = os.Getenv("ROOMS_FILE")
// 	BookingsFile = os.Getenv("BOOKINGS_FILE")
// 	UsersFile = os.Getenv("USERS_FILE")
// 	TasksFile = os.Getenv("TASKS_FILE")
// 	ServiceRequestFile = os.Getenv("SERVICE_REQUEST_FILE")
// 	FeedbackFile = os.Getenv("FEEDBACK_FILE")
// }

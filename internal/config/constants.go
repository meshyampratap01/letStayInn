package config

const (
	WelcomeMsg             = "\n            ---- Welcome to LetStayInn ----"
	UserWelcome            = "✅ Welcome, "
	AppDescription         = "Manage bookings, rooms, staff, and services with ease.\n"
	InvalidOption          = "Invalid Option. Try Again."
	SelectOption           = "\nSelect Option: "
	TitleAvailableRooms    = "\n🏨 --- Available Rooms ---"
	TitleBookRoom          = "\n🛏️ ---- Book Room ----"
	TitleCancelBooking     = "\n🗑️ ---- Cancel Booking ----"
	TitleMyBookings        = "\n📖 ---- My Bookings ----"
	TitleActiveBookings    = "📋 Your Active Bookings:"
	GuestDashboardTitle    = "\n--- Guest Dashboard ---"
	ManagerDashboardTitle  = "\n--- Manager Dashboard ---"
	EmployeeDashboardTitle = "\n--- Employee Dashboard ---"
	RoomMgmtTitle          = "\n--- Room Management ---"
	EmpMgmtTitle           = "\n--- Employee Management ---"
	FeedbackMsg            = "💬 Hello %s! We'd love to hear your thoughts."
	LoginMsg               = "\n---- 🔑 Login ---- "
	SignupMsg              = "\n---- ✍️ SignUp ---- "
)

const (
	MsgErrorFindingRooms    = "❌ Error in finding Rooms: %v"
	MsgNoAvailableRooms     = "⚠️ No available rooms."
	MsgEnterRoomNumber      = "Enter the room number to book: "
	MsgInvalidCheckInDate   = "❌ Invalid Check-in Date: %v"
	MsgInvalidCheckOutDate  = "❌ Invalid Check-out Date: %v"
	MsgBookingFailed        = "❌ Booking failed: %v"
	MsgBookingSuccess       = "✅ Room Booked Successfully!!"
	MsgFailedFetchBookings  = "❌ Failed to fetch bookings: %v"
	MsgNoBookingsToCancel   = "⚠️ No active bookings to cancel."
	MsgEnterBookingToCancel = "Enter the number of the booking to cancel: "
	MsgInvalidChoice        = "❌ Invalid choice."
	MsgCancelFailed         = "❌ Cancellation failed: %v"
	MsgCancelSuccess        = "✅ Booking cancelled successfully!"
	MsgNoBookings           = "⚠️ You have no bookings."
)

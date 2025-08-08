package handlers

import (
	"context"
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/services"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type BookingHandler struct {
	bookingService bookingService.BookingManager
	roomService    roomService.RoomServiceManager
}

func NewBookingHandler(
	bookingService bookingService.BookingManager,
	roomService roomService.RoomServiceManager,
) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		roomService:    roomService,
	}
}

func (h *BookingHandler) ViewRoomsHandler() {
	rooms, err := h.roomService.GetAvailableRooms()
	if err != nil {
		fmt.Printf("Error in finding Rooms: %v\n", err)
		services.AddBackButton()
		return
	}
	fmt.Println("\n--- Available Rooms ---")
	for _, r := range rooms {
		fmt.Printf("Room %d (%s): Rs.%.2f - %s\n", r.Number, r.Type, r.Price, r.Description)
	}
	services.AddBackButton()
}

func (h *BookingHandler) BookRoomHandler(ctx context.Context) {
	fmt.Println("\n---- Book Room ----")
	rooms, err := h.roomService.GetAvailableRooms()
	if err != nil {
		fmt.Printf("Error in finding Rooms: %v\n", err)
		services.AddBackButton()
		return
	}
	if len(rooms) == 0 {
		fmt.Println("No available rooms.")
		services.AddBackButton()
		return
	}

	for _, r := range rooms {
		fmt.Printf("Room %d (%s): Rs.%.2f - %s\n", r.Number, r.Type, r.Price, r.Description)
	}

	var roomNum int
	fmt.Print("Enter the room number to book: ")
	fmt.Scanln(&roomNum)

	var (
		checkInDateStr, checkOutDateStr string
		checkIn, checkOut               string
	)

	for {
		fmt.Print("Enter check-in date (DD-MM-YYYY): ")
		fmt.Scanln(&checkInDateStr)
		parsed, err := validators.ValidateDate(checkInDateStr)
		if err != nil {
			fmt.Println("Invalid Check-in Date:", err)
			continue
		}
		checkIn = parsed
		break
	}

	for {
		fmt.Print("Enter check-out date (DD-MM-YYYY): ")
		fmt.Scanln(&checkOutDateStr)
		parsed, err := validators.ValidateCheckoutDate(checkIn, checkOutDateStr)
		if err != nil {
			fmt.Println("Invalid Check-out Date:", err)
			continue
		}
		checkOut = parsed
		break
	}

	err = h.bookingService.BookRoom(ctx, roomNum, checkIn, checkOut)
	if err != nil {
		fmt.Printf("Booking failed: %v\n", err)
	} else {
		fmt.Println("Room Booked Successfully!!")
	}
	services.AddBackButton()
}


func (h *BookingHandler) CancelBookingHandler(ctx context.Context) {
	fmt.Println("\n---- Cancel Booking ----")

	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch bookings: %v\n", err)
		services.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		fmt.Println("No active bookings to cancel.")
		services.AddBackButton()
		return
	}

	fmt.Println("Your Active Bookings:")
	for i, b := range bookings {
		fmt.Printf("%d. Room: %d | Check-in: %s | Check-out: %s\n",
			i+1, b.RoomNum, b.CheckIn.Format("02-01-2006"), b.CheckOut.Format("02-01-2006"))
	}

	var choice int
	fmt.Print("Enter the number of the booking to cancel: ")
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(bookings) {
		fmt.Println("Invalid choice.")
		services.AddBackButton()
		return
	}

	selectedBooking := bookings[choice-1]

	err = h.bookingService.CancelBooking(ctx, selectedBooking.ID)
	if err != nil {
		fmt.Printf("Cancellation failed: %v\n", err)
	} else {
		fmt.Println("Booking cancelled successfully!")
	}
	services.AddBackButton()
}

func (h *BookingHandler) ViewMyBookingsHandler(ctx context.Context) {
	fmt.Println("\n---- My Bookings ----")
	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch bookings: %v\n", err)
		services.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		fmt.Println("You have no bookings.")
		services.AddBackButton()
		return
	}

	for _, b := range bookings {
		fmt.Printf("Room: %d | Status: %s | Check-in: %s\n",
			b.RoomNum, b.Status, b.CheckIn.Format("02-01-2006"))
	}

	services.AddBackButton()
}

package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/utils"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type BookingHandler struct {
	bookingService bookingService.IBookingService
	roomService    roomService.IRoomService
}

func NewBookingHandler(
	bookingService bookingService.IBookingService,
	roomService roomService.IRoomService,
) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		roomService:    roomService,
	}
}

func (h *BookingHandler) ViewRoomsHandler() {
	rooms, err := h.roomService.GetAvailableRooms()
	if err != nil {
		color.Red(config.MsgErrorFindingRooms, err)
		utils.AddBackButton()
		return
	}

	if len(rooms) == 0 {
		color.Yellow("No available rooms found.")
		utils.AddBackButton()
		return
	}

	color.Cyan(config.TitleAvailableRooms)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-10s %-12s %-10s %-15s %s\n",
		"Room No", "Type", "Price(Rs)", "Availability", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range rooms {
		availability := "Available"
		if !r.IsAvailable {
			availability = "Occupied"
		}

		fmt.Printf("%-10d %-12s %-10.2f %-15s %s\n",
			r.Number,
			r.Type,
			r.Price,
			availability,
			utils.TruncateString(r.Description, 30))
	}

	fmt.Println(strings.Repeat("-", 80))
	utils.AddBackButton()
}

func (h *BookingHandler) BookRoomHandler(ctx context.Context) {
	color.Cyan(config.TitleBookRoom)
	rooms, err := h.roomService.GetAvailableRooms()
	if err != nil {
		color.Red(config.MsgErrorFindingRooms, err)
		utils.AddBackButton()
		return
	}
	if len(rooms) == 0 {
		color.Yellow(config.MsgNoAvailableRooms)
		utils.AddBackButton()
		return
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-10s %-15s %-12s %-20s %-30s\n",
		"Room No", "Type", "Price (Rs)", "Availability", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range rooms {
		availability := color.GreenString("Available")
		if !r.IsAvailable {
			availability = color.RedString("Occupied")
		}

		fmt.Printf("%-10d %-15s %-12.2f %-28s %-30s\n",
			r.Number,
			r.Type,
			r.Price,
			availability,
			utils.TruncateString(r.Description, 30),
		)
	}

	fmt.Println(strings.Repeat("-", 80))

	var roomNum int
	fmt.Print(color.HiWhiteString(config.MsgEnterRoomNumber))
	fmt.Scanln(&roomNum)

	var checkInDateStr, checkOutDateStr string
	var checkIn, checkOut string

	for {
		fmt.Print(color.HiWhiteString("Enter check-in date (DD-MM-YYYY): "))
		fmt.Scanln(&checkInDateStr)
		parsed, err := validators.ValidateDate(checkInDateStr)
		if err != nil {
			color.Red(config.MsgInvalidCheckInDate, err)
			continue
		}
		checkIn = parsed
		break
	}

	for {
		fmt.Print(color.HiWhiteString("Enter check-out date (DD-MM-YYYY): "))
		fmt.Scanln(&checkOutDateStr)
		parsed, err := validators.ValidateCheckoutDate(checkIn, checkOutDateStr)
		if err != nil {
			color.Red(config.MsgInvalidCheckOutDate, err)
			continue
		}
		checkOut = parsed
		break
	}

	err = h.bookingService.BookRoom(ctx, roomNum, checkIn, checkOut)
	if err != nil {
		color.Red(config.MsgBookingFailed, err)
	} else {
		color.Green(config.MsgBookingSuccess)
	}
	utils.AddBackButton()
}

func (h *BookingHandler) CancelBookingHandler(ctx context.Context) {
	color.Cyan(config.TitleCancelBooking)

	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		color.Red(config.MsgFailedFetchBookings, err)
		utils.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		color.Yellow(config.MsgNoBookingsToCancel)
		utils.AddBackButton()
		return
	}

	color.Cyan("\nYour Active Bookings:\n")
	for i, b := range bookings {
		fmt.Println(strings.Repeat("-", 50))
		color.Yellow("%d) Room %d", i+1, b.RoomNum)
		fmt.Printf("   Check-in : %s\n", b.CheckIn.Format("02 Jan 2006"))
		fmt.Printf("   Check-out: %s\n", b.CheckOut.Format("02 Jan 2006"))
		fmt.Printf("   Status   : %s\n", b.Status)
	}

	fmt.Println(strings.Repeat("-", 50))
	var choice int
	fmt.Print(color.HiWhiteString(config.MsgEnterBookingToCancel))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(bookings) {
		color.Red(config.MsgInvalidChoice)
		utils.AddBackButton()
		return
	}

	selectedBooking := bookings[choice-1]

	fmt.Printf("\nAre you sure you want to cancel booking for Room %d (Check-in: %s, Check-out: %s)?\n",
		selectedBooking.RoomNum,
		selectedBooking.CheckIn.Format("02 Jan 2006"),
		selectedBooking.CheckOut.Format("02 Jan 2006"),
	)
	fmt.Println("1) Yes, cancel it")
	fmt.Println("2) No, keep booking")

	var confirmChoice int
	fmt.Print(color.HiWhiteString("Enter your choice: "))
	fmt.Scanln(&confirmChoice)

	if confirmChoice != 1 {
		color.Yellow("Cancellation aborted. Your booking remains active.")
		utils.AddBackButton()
		return
	}

	err = h.bookingService.CancelBooking(ctx, selectedBooking.ID)
	if err != nil {
		color.Red(config.MsgCancelFailed, err)
	} else {
		color.Green(config.MsgCancelSuccess)
	}
	utils.AddBackButton()
}

func (h *BookingHandler) ViewMyBookingsHandler(ctx context.Context) {
	fmt.Println("\n" + config.TitleMyBookings)
	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		fmt.Println("Failed to fetch bookings:", err)
		utils.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		fmt.Println("No bookings found.")
		utils.AddBackButton()
		return
	}

	for i, b := range bookings {
		roomType := "Unknown"
		room, err := h.bookingService.GetRoomByNumber(b.RoomNum)
		if err == nil && room != nil {
			roomType = string(room.Type)
		}
		fmt.Println("---------------------------------------")
		fmt.Printf(" Booking #%d\n", i+1)
		fmt.Println("---------------------------------------")
		fmt.Printf(" Room Number : %d\n", b.RoomNum)
		fmt.Printf(" Room Type   : %s\n", roomType)
		fmt.Printf(" Check-in    : %s\n", b.CheckIn.Format("02 Jan 2006, 15:04"))
		fmt.Printf(" Check-out   : %s\n", b.CheckOut.Format("02 Jan 2006, 15:04"))
		fmt.Printf(" Status      : %s\n", statusLabel(b.Status))
		fmt.Printf(" Food Req    : %s\n", utils.BoolToIcon(b.FoodReq, "Yes", "No"))
		fmt.Printf(" Cleaning    : %s\n", utils.BoolToIcon(b.CleanReq, "Yes", "No"))
		fmt.Printf(" Booked On   : %s\n", b.CreatedAt.Format("02 Jan 2006, 15:04"))
	}

	fmt.Println("---------------------------------------")
	utils.AddBackButton()
}
func statusLabel(status string) string {
	switch status {
	case models.BookingStatusBooked:
		return "Booked"
	case models.BookingStatusCancelled:
		return "Cancelled"
	case models.BookingStatusCompleted:
		return "Completed"
	default:
		return status
	}
}

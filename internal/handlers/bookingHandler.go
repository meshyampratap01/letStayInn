package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/services"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
	"github.com/meshyampratap01/letStayInn/internal/utils"
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
		color.Red(config.MsgErrorFindingRooms, err)
		services.AddBackButton()
		return
	}

	if len(rooms) == 0 {
		color.Yellow("No available rooms found.")
		services.AddBackButton()
		return
	}

	color.Cyan(config.TitleAvailableRooms)

	fmt.Printf("%-10s %-12s %-10s %-15s %s\n",
		"Room No", "Type", "Price(Rs)", "Availability", "Description")
	fmt.Println(strings.Repeat("-", 70))

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

	fmt.Println(strings.Repeat("-", 70))
	services.AddBackButton()
}

func (h *BookingHandler) BookRoomHandler(ctx context.Context) {
	color.Cyan(config.TitleBookRoom)
	rooms, err := h.roomService.GetAvailableRooms()
	if err != nil {
		color.Red(config.MsgErrorFindingRooms, err)
		services.AddBackButton()
		return
	}
	if len(rooms) == 0 {
		color.Yellow(config.MsgNoAvailableRooms)
		services.AddBackButton()
		return
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-10s %-15s %-12s %-20s %-30s\n",
		"Room No", "Type", "Price (Rs)", "Availability", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range rooms {
		availability := color.GreenString("Available")
		if !r.IsAvailable {
			availability = color.RedString("Occupied")
		}

		fmt.Printf("%-10d %-15s %-12.2f %-15s %-30s\n",
			r.Number,
			r.Type,
			r.Price,
			availability,
			utils.TruncateString(r.Description, 30),
		)
	}

	fmt.Println(strings.Repeat("=", 80))

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
	services.AddBackButton()
}

func (h *BookingHandler) CancelBookingHandler(ctx context.Context) {
	color.Cyan(config.TitleCancelBooking)

	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		color.Red(config.MsgFailedFetchBookings, err)
		services.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		color.Yellow(config.MsgNoBookingsToCancel)
		services.AddBackButton()
		return
	}

	color.Cyan(config.TitleActiveBookings)
	for i, b := range bookings {
		color.Green("%d. Room: %d | Check-in: %s | Check-out: %s",
			i+1, b.RoomNum, b.CheckIn.Format("02-01-2006"), b.CheckOut.Format("02-01-2006"))
	}

	var choice int
	fmt.Print(color.HiWhiteString(config.MsgEnterBookingToCancel))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(bookings) {
		color.Red(config.MsgInvalidChoice)
		services.AddBackButton()
		return
	}

	selectedBooking := bookings[choice-1]
	err = h.bookingService.CancelBooking(ctx, selectedBooking.ID)
	if err != nil {
		color.Red(config.MsgCancelFailed, err)
	} else {
		color.Green(config.MsgCancelSuccess)
	}
	services.AddBackButton()
}

func (h *BookingHandler) ViewMyBookingsHandler(ctx context.Context) {
	color.Cyan(config.TitleMyBookings)
	bookings, err := h.bookingService.GetUserActiveBookings(ctx)
	if err != nil {
		color.Red(config.MsgFailedFetchBookings, err)
		services.AddBackButton()
		return
	}

	if len(bookings) == 0 {
		color.Yellow(config.MsgNoBookings)
		services.AddBackButton()
		return
	}

	for _, b := range bookings {
		color.Green("Room: %d | Status: %s | Check-in: %s",
			b.RoomNum, b.Status, b.CheckIn.Format("02-01-2006"))
	}

	services.AddBackButton()
}

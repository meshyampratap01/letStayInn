package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	serviceRequest "github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

type ServiceRequestHandler struct {
	ServiceRequestService serviceRequest.IServiceRequestService
	BookingService        bookingService.IBookingService
}

func NewServiceRequestHandler(srs serviceRequest.IServiceRequestService, bs bookingService.IBookingService) *ServiceRequestHandler {
	return &ServiceRequestHandler{
		ServiceRequestService: srs,
		BookingService:        bs,
	}
}

func (s *ServiceRequestHandler) ServiceRequestHandler(ctx context.Context, reqType models.ServiceType) {
	roomNum, err := s.SelectUserRoom(ctx)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	var details string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please describe your service request in detail: ")
	details, _ = reader.ReadString('\n')
	details = strings.TrimSpace(details)

	err = s.ServiceRequestService.ServiceRequestGetter(ctx, roomNum, reqType, details)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	color.Green("Your service request has been placed successfully.")
}

func (s *ServiceRequestHandler) SelectUserRoom(ctx context.Context) (int, error) {

	userRooms, err := s.BookingService.GetUserActiveBookings(ctx)
	if err != nil {
		return -1, err
	}

	if len(userRooms) == 0 {
		return -1, fmt.Errorf("you have no active bookings")
	}

	color.Cyan("\nSelect room for the request:")
	for i, r := range userRooms {
		color.Yellow("%d. Room Num: %d", i+1, r.RoomNum)
	}

	var choice int
	fmt.Print(color.MagentaString("Enter choice: "))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(userRooms) {
		return -1, fmt.Errorf("invalid selection")
	}

	return userRooms[choice-1].RoomNum, nil
}

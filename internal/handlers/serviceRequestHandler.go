package handlers

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	serviceRequest "github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

type ServiceRequestHandler struct {
	ServiceRequestService serviceRequest.IServiceRequestService
	BookingRepo           bookingRepository.BookingRepository
}

func NewServiceRequestHandler(srs serviceRequest.IServiceRequestService, br bookingRepository.BookingRepository) *ServiceRequestHandler {
	return &ServiceRequestHandler{
		ServiceRequestService: srs,
		BookingRepo:           br,
	}
}

func (s *ServiceRequestHandler) ServiceRequestHandler(ctx context.Context, reqType models.ServiceType) {
	roomNum, err := s.SelectUserRoom(ctx)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	err = s.ServiceRequestService.ServiceRequest(ctx, roomNum, reqType)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
}

func (s *ServiceRequestHandler) SelectUserRoom(ctx context.Context) (int, error) {
	bookings, err := s.BookingRepo.GetAllBookings()
	if err != nil {
		return -1, err
	}

	var userRooms []int
	roomSet := map[string]bool{}

	for _, b := range bookings {
		if b.UserID == ctx.Value(contextkeys.UserIDKey) && b.Status != models.BookingStatusCancelled && !roomSet[b.RoomID] {
			userRooms = append(userRooms, b.RoomNum)
			roomSet[b.RoomID] = true
		}
	}

	if len(userRooms) == 0 {
		return -1, fmt.Errorf("you have no active bookings")
	}

	color.Cyan("\nSelect room for the request:")
	for i, r := range userRooms {
		color.Yellow("%d. Room Num: %d", i+1, r)
	}

	var choice int
	fmt.Print(color.MagentaString("Enter choice: "))
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(userRooms) {
		return -1, fmt.Errorf("invalid selection")
	}

	return userRooms[choice-1], nil
}

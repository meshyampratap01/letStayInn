package handlers

import (
	"context"
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/contextKeys"
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
		fmt.Println("Error:", err)
		return
	}

	err = s.ServiceRequestService.ServiceRequest(ctx, roomNum, reqType)
	if err != nil {
		fmt.Println("Error:", err)
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

	fmt.Println("\nSelect room for the request:")
	for i, r := range userRooms {
		fmt.Printf("%d. Room Num: %d\n", i+1, r)
	}

	var choice int
	fmt.Print("Enter choice: ")
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(userRooms) {
		return -1, fmt.Errorf("invalid selection")
	}

	return userRooms[choice-1], nil
}

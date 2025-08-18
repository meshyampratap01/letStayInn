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
	roomNum, err := s.SelectUserRoom(ctx,reqType)
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

func (s *ServiceRequestHandler) SelectUserRoom(ctx context.Context, reqType models.ServiceType) (int, error) {
	userRooms, err := s.BookingService.GetUserActiveBookings(ctx)
	if err != nil {
		return -1, err
	}

	if len(userRooms) == 0 {
		return -1, fmt.Errorf("you have no active bookings")
	}

	for {
		color.Cyan("\nSelect a room for your service request:\n")
		for i, r := range userRooms {
			color.Yellow("%d) Room %d", i+1, r.RoomNum)
			fmt.Printf("   Stay: %s → %s\n",
				r.CheckIn.Format("02 Jan 2006"),
				r.CheckOut.Format("02 Jan 2006"))
			fmt.Printf("   Status: %s\n", r.Status)

			if r.FoodReq {
				color.Green("   • Food service already requested\n")
			}
			if r.CleanReq {
				color.Green("   • Cleaning service already requested\n")
			}
			fmt.Println()
		}

		fmt.Print(color.HiWhiteString("Enter your choice (or 0 to cancel): "))
		var choice int
		fmt.Scanln(&choice)

		if choice == 0 {
			return -1, fmt.Errorf("operation cancelled")
		}

		if choice < 1 || choice > len(userRooms) {
			color.Red("Invalid selection. Try again.")
			continue
		}

		selected := userRooms[choice-1]

		if (reqType == models.ServiceTypeFood && selected.FoodReq) || (reqType == models.ServiceTypeCleaning && selected.CleanReq) {
			color.Yellow("\n⚠ A %s request already exists for Room %d.", reqType, selected.RoomNum)
			fmt.Print("Press ENTER to proceed anyway, or type 'r' to reselect another room: ")

			var retry string
			fmt.Scanln(&retry)

			if retry == "r" {
				continue 
			}
		}

		return selected.RoomNum, nil
	}
}



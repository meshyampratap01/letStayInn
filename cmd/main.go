package main

import (
	"fmt"
	"os"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/handlers"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/taskRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

func main() {
	userRepo := userRepository.NewFileUserRepository()
	roomRepo := roomRepository.NewRoomRepository()
	bookingRepo := bookingRepository.NewFileBookingRepository()
	feedbackRepo := feedbackRepository.NewFileFeedbackRepository()
	serviceReqRepo := serviceRequestRepository.NewFileServiceRequestRepository()
	taskRepo := taskRepository.NewFileTaskRepository()

	userService := userService.NewUserService(userRepo)
	roomService := roomService.NewRoomService(roomRepo)
	bookingService := bookingService.NewBookingService(bookingRepo, roomRepo)
	feedbackService := feedbackService.NewFeedbackService(feedbackRepo, bookingRepo)
	serviceRequestService := servicerequest.NewServiceRequestService(bookingRepo, serviceReqRepo)
	managerService := managerservice.NewManagerService(userRepo,taskRepo,serviceReqRepo)

	bookingHandler := handlers.NewBookingHandler(bookingService, roomService)
	serviceReqHandler := handlers.NewServiceRequestHandler(serviceRequestService, bookingRepo)
	managerHandler := handlers.NewManagerHandler(roomService, bookingService, userService, serviceRequestService, managerService)
	dashboardHandler := handlers.NewDashboardHandler(roomService, bookingService, feedbackService, serviceRequestService, bookingHandler, serviceReqHandler, managerHandler)
	CLIUserHandler := handlers.NewUserHandler(userService, dashboardHandler)

	for {
		fmt.Println(config.LoginMsg)
		fmt.Println("1.Signup")
		fmt.Println("2.Login")
		fmt.Println("3.Exit")
		fmt.Print("Select Option: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CLIUserHandler.SignupHandler()
		case 2:
			CLIUserHandler.LoginHandler()
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println(config.InvalidOption)
		}
	}
}

package container

import (
	"github.com/meshyampratap01/letStayInn/internal/handlers"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/taskRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

func InitHandlers() *handlers.UserHandler {
	
	userRepo := userRepository.NewFileUserRepository()
	roomRepo := roomRepository.NewRoomRepository()
	bookingRepo := bookingRepository.NewFileBookingRepository()
	feedbackRepo := feedbackRepository.NewFileFeedbackRepository()
	serviceReqRepo := serviceRequestRepository.NewFileServiceRequestRepository()
	taskRepo := taskRepository.NewFileTaskRepository()

	userSvc := userService.NewUserService(userRepo)
	roomSvc := roomService.NewRoomService(roomRepo)
	bookingSvc := bookingService.NewBookingService(bookingRepo, roomRepo, userRepo)
	feedbackSvc := feedbackService.NewFeedbackService(feedbackRepo, bookingRepo)
	serviceReqSvc := servicerequest.NewServiceRequestService(bookingRepo, serviceReqRepo)
	managerSvc := managerservice.NewManagerService(userRepo, taskRepo, serviceReqRepo, roomRepo, bookingRepo,feedbackRepo)
	employeeSvc := employeeService.NewEmployeeService(taskRepo,userRepo,roomRepo,bookingRepo,serviceReqRepo)
	

	bookingHandler := handlers.NewBookingHandler(bookingSvc, roomSvc)
	serviceReqHandler := handlers.NewServiceRequestHandler(serviceReqSvc, bookingRepo)
	managerHandler := handlers.NewManagerHandler(roomSvc, bookingSvc, userSvc, serviceReqSvc, managerSvc)
	employeeHandler := handlers.NewEmployeeHandler(employeeSvc)
	dashboardHandler := handlers.NewDashboardHandler(roomSvc, bookingSvc, feedbackSvc, serviceReqSvc, bookingHandler, serviceReqHandler, managerHandler,employeeSvc,employeeHandler)

	CLIUserHandler := handlers.NewUserHandler(userSvc, dashboardHandler,feedbackSvc)
	return CLIUserHandler
}

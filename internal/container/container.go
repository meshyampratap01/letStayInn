package container

import (
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/handlers"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/managerservice"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

type AppHandlers struct {
	UserHandler       *handlers.UserHandler
	BookingHandler    *handlers.BookingHandler
	ServiceReqHandler *handlers.ServiceRequestHandler
	ManagerHandler    *handlers.ManagerHandler
	EmployeeHandler   *handlers.EmployeeHandler
	FeedbackHandler   *handlers.FeedbackHandler
	ProfileHandler    *handlers.ProfileHandler
}

func InitHandlers() *AppHandlers {
	userRepo := userRepository.NewGormUserRepository(db.DB)
	roomRepo := roomRepository.NewGormRoomRepository(db.DB)
	bookingRepo := bookingRepository.NewGormBookingRepository(db.DB)
	feedbackRepo := feedbackRepository.NewGormFeedbackRepository(db.DB)
	serviceReqRepo := serviceRequestRepository.NewGormServiceRequestRepository(db.DB)

	userSvc := userService.NewUserService(userRepo)
	roomSvc := roomService.NewRoomService(roomRepo)
	bookingSvc := bookingService.NewBookingService(bookingRepo, roomRepo, userRepo)
	feedbackSvc := feedbackService.NewFeedbackService(feedbackRepo, bookingRepo, userRepo)
	serviceReqSvc := servicerequest.NewServiceRequestService(bookingRepo, serviceReqRepo)
	managerSvc := managerservice.NewManagerService(userRepo, serviceReqRepo, roomRepo, bookingRepo, feedbackRepo)
	employeeSvc := employeeService.NewEmployeeService(userRepo, roomRepo, bookingRepo, serviceReqRepo)

	bookingHandler := handlers.NewBookingHandler(bookingSvc, roomSvc)
	serviceReqHandler := handlers.NewServiceRequestHandler(serviceReqSvc, bookingSvc)
	managerHandler := handlers.NewManagerHandler(roomSvc, bookingSvc, userSvc, serviceReqSvc, managerSvc)
	employeeHandler := handlers.NewEmployeeHandler(employeeSvc)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackSvc)
	profileHandler := handlers.NewProfileHandler(userSvc)


	userHandler := handlers.NewUserHandler(userSvc)

	return &AppHandlers{
		UserHandler:       userHandler,
		BookingHandler:    bookingHandler,
		ServiceReqHandler: serviceReqHandler,
		ManagerHandler:    managerHandler,
		EmployeeHandler:   employeeHandler,
		FeedbackHandler:   feedbackHandler,
		ProfileHandler:    profileHandler,
	}
}

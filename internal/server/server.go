package server

import (
	"log"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/container"
)


func withAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.JWTAuthMiddleware(h).ServeHTTP(w, r)
	}
}

func withCorsHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.CorsMiddleWare(h).ServeHTTP(w, r)
	}
}


func registerUserRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"signup", withCorsHeader(h.UserHandler.SignupHTTPHandler))
	mux.HandleFunc("OPTIONS "+prefix+"login", withCorsHeader(h.UserHandler.LoginHTTPHandler))
	mux.HandleFunc("POST "+prefix+"login", withCorsHeader(h.UserHandler.LoginHTTPHandler))
}

func registerRoomRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"rooms", withCorsHeader(withAuth(h.BookingHandler.GetRoomsByRoleHTTP)))// Use ?available=true for manager
	mux.HandleFunc("OPTIONS "+prefix+"rooms", withCorsHeader(withAuth(h.BookingHandler.GetRoomsByRoleHTTP)))
	mux.HandleFunc("POST "+prefix+"rooms", withCorsHeader(withAuth(h.ManagerHandler.AddRoomHTTP)))
	mux.HandleFunc("PUT "+prefix+"rooms/{roomNum}", withCorsHeader(withAuth(h.ManagerHandler.UpdateRoomHTTP)))
	mux.HandleFunc("DELETE "+prefix+"rooms/{roomNum}", withCorsHeader(withAuth(h.ManagerHandler.DeleteRoomHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"rooms/{roomNum}", withCorsHeader(withAuth(h.ManagerHandler.DeleteRoomHTTP)))
}

func registerBookingRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"bookings", withCorsHeader(withAuth(h.BookingHandler.GetBookingsByRoleHTTP)))// Use ?all=true for manager
	mux.HandleFunc("OPTIONS "+prefix+"bookings", withCorsHeader(withAuth(h.BookingHandler.GetBookingsByRoleHTTP)))
	mux.HandleFunc("POST "+prefix+"bookings", withAuth(h.BookingHandler.BookRoomHTTP))
	mux.HandleFunc("DELETE "+prefix+"bookings/{bookingId}", withCorsHeader(withAuth(h.BookingHandler.CancelBookingHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"bookings/{bookingId}", withCorsHeader(withAuth(h.BookingHandler.CancelBookingHTTP)))
}

func registerFeedbackRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"feedbacks", withAuth(h.FeedbackHandler.SubmitFeedbackHTTP))
	mux.HandleFunc("GET "+prefix+"feedbacks", withCorsHeader(withAuth(h.ManagerHandler.ListAllFeedbackHTTP))) 
	mux.HandleFunc("OPTIONS "+prefix+"feedbacks", withCorsHeader(withAuth(h.ManagerHandler.ListAllFeedbackHTTP))) 
}

func registerEmployeeRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"employees", withCorsHeader(withAuth(h.ManagerHandler.ListEmployeesHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"employees", withCorsHeader(withAuth(h.ManagerHandler.ListEmployeesHTTP)))
	mux.HandleFunc("POST "+prefix+"employees", withAuth(h.ManagerHandler.CreateEmployeeHTTP))
	mux.HandleFunc("DELETE "+prefix+"employees/{employeeId}", withCorsHeader(withAuth(h.ManagerHandler.DeleteEmployeeHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"employees/{employeeId}", withCorsHeader(withAuth(h.ManagerHandler.DeleteEmployeeHTTP)))
	mux.HandleFunc("PUT "+prefix+"employees/{employeeEmail}/availability", withCorsHeader(withAuth(h.ManagerHandler.UpdateEmployeeAvailabilityHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"employees/{employeeEmail}/availability", withCorsHeader(withAuth(h.ManagerHandler.UpdateEmployeeAvailabilityHTTP)))
	mux.HandleFunc("GET "+prefix+"employee/service-requests", withAuth(h.EmployeeHandler.ViewAssignedServiceRequestsHTTP))
	mux.HandleFunc("PUT "+prefix+"employee/service-requests/{serviceRequestId}/status", withAuth(h.EmployeeHandler.UpdateServiceRequestStatusHTTP))
	mux.HandleFunc("PUT "+prefix+"employee/availability", withAuth(h.EmployeeHandler.ToggleAvailabilityHTTP))
}

func registerServiceRequestRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"service-requests", withCorsHeader(withAuth(h.ServiceReqHandler.SubmitServiceRequestHTTP)))
	mux.HandleFunc("GET "+prefix+"service-requests", withCorsHeader(withAuth(h.ManagerHandler.ListUnassignedServiceRequestsHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"service-requests", withCorsHeader(withAuth(h.ManagerHandler.ListUnassignedServiceRequestsHTTP)))
	mux.HandleFunc("POST "+prefix+"service-requests/{requestId}/assign", withCorsHeader(withAuth(h.ManagerHandler.AssignServiceRequestHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"service-requests/{requestId}/assign", withCorsHeader(withAuth(h.ManagerHandler.AssignServiceRequestHTTP)))
	mux.HandleFunc("PUT "+prefix+"service-requests/{requestId}/status", withAuth(h.ManagerHandler.UpdateServiceRequestStatusHTTP))
	mux.HandleFunc("DELETE "+prefix+"service-requests/{requestId}", withAuth(h.ManagerHandler.CancelServiceRequestHTTP))
}

func registerProfileRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"profile", withCorsHeader(withAuth(h.ProfileHandler.ViewProfileHTTP)))
	mux.HandleFunc("PUT "+prefix+"profile", withCorsHeader(withAuth(h.ProfileHandler.UpdateProfileHTTP)))
	mux.HandleFunc("OPTIONS "+prefix+"profile", withCorsHeader(withAuth(h.ProfileHandler.UpdateProfileHTTP)))
}


func NewRouter(handlers *container.AppHandlers) http.Handler {
	mux := http.NewServeMux()
	apiPrefix := "/api/v1/"

	registerUserRoutes(mux, apiPrefix, handlers)
	registerRoomRoutes(mux, apiPrefix, handlers)
	registerBookingRoutes(mux, apiPrefix, handlers)
	registerFeedbackRoutes(mux, apiPrefix, handlers)
	registerEmployeeRoutes(mux, apiPrefix, handlers)
	registerServiceRequestRoutes(mux, apiPrefix, handlers)
	registerProfileRoutes(mux, apiPrefix, handlers)

	return mux
}

func StartServer() {
	handlers := container.InitHandlers()
	router := NewRouter(handlers)
	log.Println("HTTP API server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

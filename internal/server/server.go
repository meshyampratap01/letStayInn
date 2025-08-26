package server

import (
	"log"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/container"
)

// helper to apply JWTAuth middleware
func withAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.JWTAuthMiddleware(h).ServeHTTP(w, r)
	}
}

// --- Route Groups ---
func registerUserRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"signup", h.UserHandler.SignupHTTPHandler)
	mux.HandleFunc("POST "+prefix+"login", h.UserHandler.LoginHTTPHandler)
}

func registerRoomRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"rooms", withAuth(h.BookingHandler.GetRoomsByRoleHTTP))// Use ?available=true for manager
	mux.HandleFunc("POST "+prefix+"rooms", withAuth(h.ManagerHandler.AddRoomHTTP))
	mux.HandleFunc("PUT "+prefix+"rooms/{roomNum}", withAuth(h.ManagerHandler.UpdateRoomHTTP))
	mux.HandleFunc("DELETE "+prefix+"rooms/{roomNum}", withAuth(h.ManagerHandler.DeleteRoomHTTP))
}

func registerBookingRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"bookings", withAuth(h.BookingHandler.GetBookingsByRoleHTTP)) // Use ?all=true for manager
	mux.HandleFunc("POST "+prefix+"bookings", withAuth(h.BookingHandler.BookRoomHTTP))
	mux.HandleFunc("DELETE "+prefix+"bookings/{bookingId}", withAuth(h.BookingHandler.CancelBookingHTTP))
}

func registerFeedbackRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"feedback", withAuth(h.FeedbackHandler.SubmitFeedbackHTTP))
	mux.HandleFunc("GET "+prefix+"feedback", withAuth(h.ManagerHandler.ListAllFeedbackHTTP)) // Use ?all=true for manager
}

func registerEmployeeRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"employees", withAuth(h.ManagerHandler.ListEmployeesHTTP))
	mux.HandleFunc("POST "+prefix+"employees", withAuth(h.ManagerHandler.CreateEmployeeHTTP))
	mux.HandleFunc("DELETE "+prefix+"employees/{employeeId}", withAuth(h.ManagerHandler.DeleteEmployeeHTTP))
	mux.HandleFunc("PUT "+prefix+"employees/{employeeId}/availability", withAuth(h.ManagerHandler.UpdateEmployeeAvailabilityHTTP))
	mux.HandleFunc("GET "+prefix+"employee/service-requests", withAuth(h.EmployeeHandler.ViewAssignedServiceRequestsHTTP))
	mux.HandleFunc("PUT "+prefix+"employee/service-requests/{requestId}/status", withAuth(h.EmployeeHandler.UpdateServiceRequestStatusHTTP))
	mux.HandleFunc("PUT "+prefix+"employee/availability", withAuth(h.EmployeeHandler.ToggleAvailabilityHTTP))
}

func registerServiceRequestRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("POST "+prefix+"service-requests", withAuth(h.ServiceReqHandler.SubmitServiceRequestHTTP))
	mux.HandleFunc("GET "+prefix+"service-requests", withAuth(h.ManagerHandler.ListUnassignedServiceRequestsHTTP))
	mux.HandleFunc("POST "+prefix+"service-requests/{requestId}/assign", withAuth(h.ManagerHandler.AssignServiceRequestHTTP))
	mux.HandleFunc("PUT "+prefix+"service-requests/{requestId}/status", withAuth(h.ManagerHandler.UpdateServiceRequestStatusHTTP))
	mux.HandleFunc("DELETE "+prefix+"service-requests/{requestId}", withAuth(h.ManagerHandler.CancelServiceRequestHTTP))
}

func registerProfileRoutes(mux *http.ServeMux, prefix string, h *container.AppHandlers) {
	mux.HandleFunc("GET "+prefix+"profile", withAuth(h.ProfileHandler.ViewProfileHTTP))
	mux.HandleFunc("PUT "+prefix+"profile", withAuth(h.ProfileHandler.UpdateProfileHTTP))
}

// --- MAIN ROUTER ---
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

// StartServer starts the HTTP server
func StartServer() {
	handlers := container.InitHandlers()
	router := NewRouter(handlers)
	log.Println("HTTP API server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

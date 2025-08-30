package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/response"

	"go.uber.org/zap"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type UserHandler struct {
	userService userService.IUserService
}

func NewUserHandler(us userService.IUserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// SignupRequest represents the expected JSON body for signup
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the expected JSON body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupHTTPHandler handles user signup via HTTP (JSON)
func (u *UserHandler) SignupHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid signup request body", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if err := validators.ValidateEmail(req.Email); err != nil {
		logger.Log.Warn("Invalid email during signup", zap.String("email", req.Email))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid email")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if err := validators.ValidatePassword(req.Password); err != nil {
		logger.Log.Warn("Invalid password during signup", zap.String("email", req.Email))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid password")
		json.NewEncoder(w).Encode(resp)
		return
	}

	msg, err := u.userService.Signup(req.Name, req.Email, req.Password, int(models.RoleGuest))
	if err != nil {
		logger.Log.Error("Signup failed", zap.Error(err), zap.String("email", req.Email))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("User signed up successfully", zap.String("email", req.Email))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := response.NewSuccessResponse(http.StatusCreated, msg, nil)
	json.NewEncoder(w).Encode(resp)
}

// LoginHTTPHandler handles user login via HTTP (JSON)
func (u *UserHandler) LoginHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid login request body", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	user, err := u.userService.Login(req.Email, req.Password)
	if err != nil {
		logger.Log.Warn("Login failed", zap.Error(err), zap.String("email", req.Email))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	token, err := auth.GenerateJWT(user.ID, user.Name, user.Role.String(), config.JWTExpirationMinute)
	if err != nil {
		logger.Log.Error("Failed to generate JWT during login", zap.Error(err), zap.String("userID", user.ID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, "Failed to generate JWT")
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("User logged in successfully", zap.String("userID", user.ID), zap.String("email", req.Email))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewSuccessResponse(http.StatusOK, "Login successful", map[string]string{"token": token}))
}

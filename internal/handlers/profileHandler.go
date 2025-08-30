package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/response"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	userService userService.IUserService
}

func NewProfileHandler(userService userService.IUserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// GET /api/v1/profile
func (h *ProfileHandler) ViewProfileHTTP(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		logger.Log.Error("User ID not found in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "User ID not found in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Log.Error("Error fetching user profile", zap.String("userID", userID), zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, "Error fetching user profile")
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("User profile fetched", zap.String("userID", userID))
	resp := dto.UserProfileDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.String(),
		Available: user.Available,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.NewSuccessResponse(http.StatusOK, "User profile fetched successfully", resp))
}

// PUT /api/v1/profile
func (h *ProfileHandler) UpdateProfileHTTP(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		logger.Log.Error("User ID not found in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "User ID not found in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Log.Error("Error fetching user profile", zap.String("userID", userID), zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, "Error fetching user profile")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		Name     *string `json:"name,omitempty"`
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for profile update", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	updated := false
	if req.Name != nil && len(*req.Name) >= 2 {
		user.Name = *req.Name
		updated = true
	}
	if req.Email != nil && len(*req.Email) >= 5 && containsAt(*req.Email) {
		user.Email = *req.Email
		updated = true
	}
	if req.Password != nil && len(*req.Password) >= 6 {
		user.Password = auth.HashPassword(*req.Password)
		updated = true
	}
	if !updated {
		logger.Log.Warn("No valid fields to update in profile", zap.String("userID", userID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "No valid fields to update")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = h.userService.UpdateUser(user)
	if err != nil {
		logger.Log.Error("Error updating profile", zap.String("userID", userID), zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, "Error updating profile")
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Profile updated successfully", zap.String("userID", userID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.NewSuccessResponse(http.StatusOK, "Profile updated successfully", nil)
	json.NewEncoder(w).Encode(resp)
}

func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}

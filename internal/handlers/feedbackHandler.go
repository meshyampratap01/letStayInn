package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meshyampratap01/letStayInn/internal/response"

	"go.uber.org/zap"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
)

type FeedbackHandler struct {
	feedbackService feedbackService.IFeedbackService
}

func NewFeedbackHandler(feedbackService feedbackService.IFeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

// POST /api/v1/feedback
func (fh *FeedbackHandler) SubmitFeedbackHTTP(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid role in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "Unauthorized: invalid role in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if role != "Guest" {
		logger.Log.Warn("Forbidden: guest access required", zap.String("role", role))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		resp := response.NewErrorResponse(http.StatusForbidden, "Forbidden: guest access required")
		json.NewEncoder(w).Encode(resp)
		return
	}
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid user ID in context", zap.Any("context", r.Context()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		resp := response.NewErrorResponse(http.StatusUnauthorized, "Unauthorized: invalid user ID in context")
		json.NewEncoder(w).Encode(resp)
		return
	}
	var req struct {
		Message string `json:"message"`
		Rating  int    `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for feedback submission", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewErrorResponse(http.StatusBadRequest, "Invalid request body")
		json.NewEncoder(w).Encode(resp)
		return
	}
	err := fh.feedbackService.SubmitFeedback(r.Context(), req.Message, req.Rating)
	if err != nil {
		logger.Log.Error("Failed to submit feedback", zap.Error(err), zap.String("userID", userID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewErrorResponse(http.StatusInternalServerError, err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Log.Info("Feedback submitted", zap.String("userID", userID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := response.NewSuccessResponse(http.StatusCreated, "Feedback submitted", nil)
	json.NewEncoder(w).Encode(resp)
}

package handlers

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, "Unauthorized: invalid role in context", http.StatusUnauthorized)
		return
	}
	if role != "Guest" {
		logger.Log.Warn("Forbidden: guest access required", zap.String("role", role))
		http.Error(w, "Forbidden: guest access required", http.StatusForbidden)
		return
	}
	userIDVal := r.Context().Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok {
		logger.Log.Error("Unauthorized: invalid user ID in context", zap.Any("context", r.Context()))
		http.Error(w, "Unauthorized: invalid user ID in context", http.StatusUnauthorized)
		return
	}
	var req struct {
		Message string `json:"message"`
		Rating  int    `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body for feedback submission", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := fh.feedbackService.SubmitFeedback(r.Context(), userID, req.Rating)
	if err != nil {
		logger.Log.Error("Failed to submit feedback", zap.Error(err), zap.String("userID", userID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Feedback submitted", zap.String("userID", userID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Feedback submitted"})
}

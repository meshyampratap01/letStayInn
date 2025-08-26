package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
)

type FeedbackHandler struct {
	feedbackService feedbackService.IFeedbackService
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

func NewFeedbackHandler(feedbackService feedbackService.IFeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

func (h *FeedbackHandler) SubmitFeedback(ctx context.Context) error {
	if name, ok := ctx.Value(contextkeys.UserNameKey).(string); ok {
		color.Cyan(config.FeedbackMsg, name)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Submit Feedback ---")
	fmt.Print("Enter your feedback: ")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("feedback cannot be empty")
	}

	fmt.Print("Rate your experience (1-5, optional, press Enter to skip): ")
	ratingInput, _ := reader.ReadString('\n')
	ratingInput = strings.TrimSpace(ratingInput)

	rating := 0
	if ratingInput != "" {
		if r, err := strconv.Atoi(ratingInput); err == nil && r >= 1 && r <= 5 {
			rating = r
		} else {
			color.Yellow("Invalid rating, skipping...")
		}
	}

	if err := h.feedbackService.SubmitFeedback(ctx, message, rating); err != nil {
		return fmt.Errorf("error submitting feedback: %w", err)
	}

	color.Green("✅ Thank you for your feedback!")
	return nil
}

// func (fh *FeedbackHandler) GetFeedbacksHTTP(w http.ResponseWriter, r *http.Request) {
// 	// Not implemented: no method in IFeedbackService for fetching by user/role
// }

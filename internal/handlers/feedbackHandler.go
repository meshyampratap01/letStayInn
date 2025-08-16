package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
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

package feedbackService

import (
	"context"

	"github.com/meshyampratap01/letStayInn/internal/models"
)

type IFeedbackService interface {
	SubmitFeedback(ctx context.Context, message string, rating int) error
	DeleteFeedback(feedbackID string) error
	ViewAllFeedback() ([]models.Feedback, error)
}

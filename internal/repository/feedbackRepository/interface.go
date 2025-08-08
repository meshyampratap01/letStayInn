package feedbackRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type FeedbackRepository interface {
	SaveFeedback(models.Feedback) error
	GetAllFeedback() ([]models.Feedback, error)
}
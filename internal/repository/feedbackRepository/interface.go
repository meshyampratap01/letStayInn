//go:generate mockgen -source=interface.go -destination=../../mocks/mock_feedbackRepository.go -package=mocks
package feedbackRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type FeedbackRepository interface {
	SaveFeedback(models.Feedback) error
	GetAllFeedback() ([]models.Feedback, error)
	DeleteFeedback(string) error
}

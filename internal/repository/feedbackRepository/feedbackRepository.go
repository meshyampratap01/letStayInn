package feedbackRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileFeedbackRepository struct{}

func NewFileFeedbackRepository() FeedbackRepository{
	return &FileFeedbackRepository{}
}

func (db *FileFeedbackRepository)SaveFeedback(f models.Feedback) error {
	var feedbacks []models.Feedback
	if err := storage.ReadJson(config.FeedbackFile, &feedbacks); err != nil {
		return err
	}
	feedbacks = append(feedbacks, f)
	return storage.WriteJson(config.FeedbackFile, feedbacks)
}


func (repo *FileFeedbackRepository) GetAllFeedback() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := storage.ReadJson(config.FeedbackFile, &feedbacks)
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}
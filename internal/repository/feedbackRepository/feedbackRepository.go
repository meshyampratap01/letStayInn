package feedbackRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/gorm"
)

type GormFeedbackRepository struct {
	db *gorm.DB
}

func NewGormFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &GormFeedbackRepository{db: db}
}

func (r *GormFeedbackRepository) SaveFeedback(f models.Feedback) error {
	return r.db.Create(&f).Error
}

func (r *GormFeedbackRepository) GetAllFeedback() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.db.Find(&feedbacks).Error
	return feedbacks, err
}


func (r *GormFeedbackRepository) DeleteFeedback(feedbackID string) error {
    err := r.db.Delete(&models.Feedback{}, "id = ?", feedbackID).Error
    return err
}

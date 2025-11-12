package feedbackService

import "context"

//go:generate mockgen -source=interface.go -destination=../../mocks/mock_feedbackService.go -package=mocks

type IFeedbackService interface {
	SubmitFeedback(ctx context.Context, message string, rating int) error
	DeleteFeedback(feedbackID string) error
}

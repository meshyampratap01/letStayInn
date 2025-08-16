package feedbackService

import "context"


type IFeedbackService interface {
	SubmitFeedback(ctx context.Context, message string, rating int) error
}

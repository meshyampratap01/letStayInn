package feedbackService

import "context"


type IFeedbackService interface {
	SubmitFeedback(ctx context.Context) error
}

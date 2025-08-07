package feedbackService

import "context"


type FeedbackServiceManager interface {
	SubmitFeedback(ctx context.Context) error
}

package feedbackService

import "context"

// FeedbackService defines the contract for feedback-related operations.
type FeedbackServiceManager interface {
	SubmitFeedback(ctx context.Context) error
}

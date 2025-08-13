package feedbackService

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type FeedbackService struct {
	feedbackRepo feedbackRepository.FeedbackRepository
	bookingRepo  bookingRepository.BookingRepository
}

func NewFeedbackService(feedbackRepo feedbackRepository.FeedbackRepository,
	bookingRepo bookingRepository.BookingRepository) IFeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
		bookingRepo:  bookingRepo,
	}
}

func (s *FeedbackService) SubmitFeedback(ctx context.Context) error {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("invalid or missing user ID in context")
	}
	bookings, err := s.bookingRepo.GetBookingsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch bookings: %w", err)
	}

	var eligibleBooking *models.Booking
	for i := range bookings {
		if bookings[i].Status == models.BookingStatusCompleted || bookings[i].Status == models.BookingStatusBooked {
			eligibleBooking = &bookings[i]
			break
		}
	}

	if eligibleBooking == nil {
		return errors.New("you need at least one completed or active booking to submit feedback")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--- Submit Feedback ---")
	fmt.Print("Enter your feedback: ")

	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	if message == "" {
		return errors.New("feedback cannot be empty")
	}

	feedback := models.Feedback{
		ID:        utils.NewUUID(),
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.feedbackRepo.SaveFeedback(feedback); err != nil {
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	eligibleBooking.FeedbackID = append(eligibleBooking.FeedbackID, feedback.ID)
	if err := s.bookingRepo.UpdateBooking(*eligibleBooking); err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	fmt.Println("Thank you for the feedback!")
	return nil
}

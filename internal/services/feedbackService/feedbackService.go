package feedbackService

import (
	"context"
	"errors"
	"fmt"
	"time"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type FeedbackService struct {
	feedbackRepo feedbackRepository.FeedbackRepository
	bookingRepo  bookingRepository.BookingRepository
	userRepo     userRepository.UserRepository
}

func NewFeedbackService(feedbackRepo feedbackRepository.FeedbackRepository,
	bookingRepo bookingRepository.BookingRepository,
	userRepo userRepository.UserRepository) IFeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
		bookingRepo:  bookingRepo,
		userRepo:     userRepo,
	}
}

func (s *FeedbackService) SubmitFeedback(ctx context.Context, message string, rating int) error {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("invalid or missing user ID in context")
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

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	feedback := models.Feedback{
		ID:        utils.NewUUID(),
		UserID:    userID,
		UserName:  user.Name,
		Message:   message,
		CreatedAt: time.Now(),
		BookingID: eligibleBooking.ID,
		RoomNum:   eligibleBooking.RoomNum,
		Rating:    rating,
	}

	if err := s.feedbackRepo.SaveFeedback(feedback); err != nil {
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	eligibleBooking.FeedbackID = append(eligibleBooking.FeedbackID, feedback.ID)
	if err := s.bookingRepo.UpdateBooking(*eligibleBooking); err != nil {
		return fmt.Errorf("failed to update booking with feedback: %w", err)
	}

	return nil
}

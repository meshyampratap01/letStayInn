package feedbackService

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	gomock "go.uber.org/mock/gomock"
// 	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
// 	"github.com/meshyampratap01/letStayInn/internal/models"
// 	mockBooking "github.com/meshyampratap01/letStayInn/internal/mocks" 
// )

// func TestSubmitFeedback(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBookingRepo := mockBooking.NewMockBookingRepository(ctrl)
// 	mockFeedbackRepo := mockBooking.NewMockFeedbackRepository(ctrl)
// 	mockUserRepo := mockBooking.NewMockUserRepository(ctrl)

// 	service :=  NewFeedbackService(
// 		mockFeedbackRepo,
// 		mockBookingRepo,
// 		mockUserRepo,
// 	)

// 	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, "user-123")

// 	tests := []struct {
// 		name          string
// 		ctx           context.Context
// 		setupMocks    func()
// 		expectedError bool
// 	}{
// 		{
// 			name: "missing user id in context",
// 			ctx:  context.Background(),
// 			setupMocks: func() {
// 				// no mocks
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "booking repo error",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return(nil, errors.New("db error"))
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "no eligible booking",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return([]models.Booking{
// 						{ID: "b1", Status: models.BookingStatusCancelled},
// 					}, nil)
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "user repo error",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return([]models.Booking{
// 						{ID: "b1", Status: models.BookingStatusCompleted, RoomNum: 101},
// 					}, nil)

// 				mockUserRepo.EXPECT().
// 					GetUserByID("user-123").
// 					Return(nil, errors.New("user not found"))
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "feedback repo error",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return([]models.Booking{
// 						{ID: "b1", Status: models.BookingStatusCompleted, RoomNum: 101},
// 					}, nil)

// 				mockUserRepo.EXPECT().
// 					GetUserByID("user-123").
// 					Return(&models.User{ID: "user-123", Name: "Shyam"}, nil)

// 				mockFeedbackRepo.EXPECT().
// 					SaveFeedback(gomock.Any()).
// 					Return(errors.New("save error"))
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "update booking error",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return([]models.Booking{
// 						{ID: "b1", Status: models.BookingStatusCompleted, RoomNum: 101},
// 					}, nil)

// 				mockUserRepo.EXPECT().
// 					GetUserByID("user-123").
// 					Return(&models.User{ID: "user-123", Name: "Shyam"}, nil)

// 				mockFeedbackRepo.EXPECT().
// 					SaveFeedback(gomock.Any()).
// 					Return(nil)

// 				mockBookingRepo.EXPECT().
// 					UpdateBooking(gomock.Any()).
// 					Return(errors.New("update error"))
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name: "success",
// 			ctx:  ctx,
// 			setupMocks: func() {
// 				mockBookingRepo.EXPECT().
// 					GetBookingsByUserID("user-123").
// 					Return([]models.Booking{
// 						{ID: "b1", Status: models.BookingStatusBooked, RoomNum: 101},
// 					}, nil)

// 				mockUserRepo.EXPECT().
// 					GetUserByID("user-123").
// 					Return(&models.User{ID: "user-123", Name: "Shyam"}, nil)

// 				mockFeedbackRepo.EXPECT().
// 					SaveFeedback(gomock.Any()).
// 					Return(nil)

// 				mockBookingRepo.EXPECT().
// 					UpdateBooking(gomock.Any()).
// 					Return(nil)
// 			},
// 			expectedError: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()
// 			err := service.SubmitFeedback(tt.ctx, "Great stay", 5)
// 			if tt.expectedError && err == nil {
// 				t.Errorf("expected error but got nil")
// 			}
// 			if !tt.expectedError && err != nil {
// 				t.Errorf("did not expect error but got %v", err)
// 			}
// 		})
// 	}
// }

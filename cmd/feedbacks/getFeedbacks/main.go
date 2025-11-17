package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	authmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/authMiddleware"
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/feedbackRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
)

var feedbackSvc feedbackService.IFeedbackService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	feedbackRepo := feedbackRepository.NewFeedbackRepo(dynamoDB, tableName)
	feedbackSvc = feedbackService.NewFeedbackService(feedbackRepo, bookingRepository.NewBookingRepo(dynamoDB, tableName), userRepository.NewUserRepo(dynamoDB, tableName))
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	feedbacks, err := feedbackSvc.ViewAllFeedback()
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	dtos := make([]dto.FeedbackDTO, 0, len(feedbacks))
	for _, f := range feedbacks {
		dtos = append(dtos, dto.FeedbackDTO{
			ID:        f.ID,
			UserID:    f.UserID,
			UserName:  f.UserName,
			Message:   f.Message,
			CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
			RoomNum:   f.RoomNum,
			BookingID: f.BookingID,
			Rating:    f.Rating,
		})
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Feedback fetched successfully", dtos)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

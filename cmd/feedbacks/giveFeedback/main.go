package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/db"
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
	bookingRepo := bookingRepository.NewBookingRepo(dynamoDB, tableName)
	userRepo := userRepository.NewUserRepo(dynamoDB, tableName)
	feedbackSvc = feedbackService.NewFeedbackService(feedbackRepo, bookingRepo, userRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusUnauthorized, "Unauthorized: invalid role in context"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusUnauthorized,
			Body:       string(body),
		}, nil
	}

	if role != "Guest" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Forbidden: guest access required"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
			Body:       string(body),
		}, nil
	}

	var feedbackReq struct {
		Message string `json:"message"`
		Rating  int    `json:"rating"`
	}

	if err := json.Unmarshal([]byte(req.Body), &feedbackReq); err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid request body"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	err := feedbackSvc.SubmitFeedback(ctx, feedbackReq.Message, feedbackReq.Rating)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusCreated, "Feedback submitted", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusCreated,
		Body:       string(body),
	}, nil
}

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
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/bookingService"
)

var bookingSvc bookingService.IBookingService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	bookingRepo := bookingRepository.NewBookingRepo(dynamoDB, tableName)
	roomRepo := roomRepository.NewRoomRepo(dynamoDB, tableName)
	userRepo := userRepository.NewUserRepo(dynamoDB, tableName)
	bookingSvc = bookingService.NewBookingService(bookingRepo, roomRepo, userRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var bookingReq struct {
		RoomNumber   int    `json:"room_number"`
		CheckInDate  string `json:"check_in_date"`
		CheckOutDate string `json:"check_out_date"`
	}

	if err := json.Unmarshal([]byte(req.Body), &bookingReq); err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid request body"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	userIDVal := ctx.Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusUnauthorized, "User ID not found in context"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusUnauthorized,
			Body:       string(body),
		}, nil
	}

	err := bookingSvc.BookRoom(ctx, bookingReq.RoomNumber, bookingReq.CheckInDate, bookingReq.CheckOutDate)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusCreated, "Room booked successfully", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusCreated,
		Body:       string(body),
	}, nil
}

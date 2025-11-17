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
	"github.com/meshyampratap01/letStayInn/internal/models"
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
	bookings, err := bookingSvc.GetUserActiveBookings(ctx)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, "Failed to fetch bookings"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	bookingDTOs := convertModelsToBookingDTOs(bookings)
	resp := response.NewSuccessResponse(constants.StatusOK, "Bookings fetched successfully", bookingDTOs)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

func convertModelsToBookingDTOs(bookings []models.Booking) []dto.BookingDTO {
	bookingDTOs := make([]dto.BookingDTO, 0, len(bookings))
	for _, booking := range bookings {
		bookingDTOs = append(bookingDTOs, dto.BookingDTO{
			ID:         booking.ID,
			RoomNumber: booking.RoomNum,
			Status:     booking.Status,
			FoodReq:    booking.FoodReq,
			CleanReq:   booking.CleanReq,
			CheckIn:    booking.CheckIn.Format("2006-01-02"),
			CheckOut:   booking.CheckOut.Format("2006-01-02"),
		})
	}
	return bookingDTOs
}

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
	"github.com/meshyampratap01/letStayInn/internal/dto"
	authmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/authMiddleware"
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
)

var employeeSvc employeeService.IEmployeeService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	userRepo := userRepository.NewUserRepo(dynamoDB, tableName)
	roomRepo := roomRepository.NewRoomRepo(dynamoDB, tableName)
	bookingRepo := bookingRepository.NewBookingRepo(dynamoDB, tableName)
	serviceReqRepo := serviceRequestRepository.NewServiceRequestRepo(dynamoDB, tableName)
	employeeSvc = employeeService.NewEmployeeService(userRepo, roomRepo, bookingRepo, serviceReqRepo)
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

	userIDVal := ctx.Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusUnauthorized, "Unauthorized: invalid user ID in context"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusUnauthorized,
			Body:       string(body),
		}, nil
	}

	if role != "KitchenStaff" && role != "CleaningStaff" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Forbidden: employee access required"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
			Body:       string(body),
		}, nil
	}

	requests, err := employeeSvc.GetAssignedServiceRequests(userID)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	resp := make([]dto.ServiceRequestDTO, 0, len(requests))
	for _, sr := range requests {
		resp = append(resp, dto.ServiceRequestDTO{
			ID:         sr.ID,
			RoomNum:    sr.RoomNum,
			Type:       string(sr.Type),
			Details:    sr.Details,
			Status:     string(sr.Status),
			EmployeeID: sr.AssignedTo,
			IsAssigned: sr.IsAssigned,
		})
	}

	body, _ := json.Marshal(response.NewSuccessResponse(constants.StatusOK, "Assigned service requests fetched successfully", resp))
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

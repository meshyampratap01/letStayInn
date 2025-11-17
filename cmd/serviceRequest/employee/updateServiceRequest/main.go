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
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/bookingRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/employeeService"
)

var empSvc employeeService.IEmployeeService

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
	empSvc = employeeService.NewEmployeeService(userRepo, roomRepo, bookingRepo, serviceReqRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok || role == "" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusUnauthorized, "Unauthorized: invalid role in context"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusUnauthorized,
			Body:       string(body),
		}, nil
	}

	userIDVal := ctx.Value(contextkeys.UserIDKey)
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
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

	requestID := req.PathParameters["serviceRequestId"]
	if requestID == "" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Missing service request id"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	var statusReq struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(req.Body), &statusReq); err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid request body"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	var newStatus models.ServiceStatus
	switch statusReq.Status {
	case string(models.ServiceStatusPending), string(models.ServiceStatusInProgress), string(models.ServiceStatusDone):
		newStatus = models.ServiceStatus(statusReq.Status)
	default:
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid status"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	assignedRequests, err := empSvc.GetAssignedServiceRequests(userID)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	found := false
	for _, sr := range assignedRequests {
		if sr.ID == requestID {
			found = true
			break
		}
	}

	if !found {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Service request not assigned to you"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
			Body:       string(body),
		}, nil
	}

	err = empSvc.UpdateServiceRequestStatus(requestID, newStatus)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Service request status updated", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

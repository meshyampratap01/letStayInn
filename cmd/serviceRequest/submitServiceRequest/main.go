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
	"github.com/meshyampratap01/letStayInn/internal/repository/serviceRequestRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/servicerequest"
)

var serviceReqSvc servicerequest.IServiceRequestService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	bookingRepo := bookingRepository.NewBookingRepo(dynamoDB, tableName)
	serviceReqRepo := serviceRequestRepository.NewServiceRequestRepo(dynamoDB, tableName)
	serviceReqSvc = servicerequest.NewServiceRequestService(bookingRepo, serviceReqRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok || role != "Guest" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Forbidden: guest access required"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
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

	var serviceReq struct {
		RoomNum int    `json:"room_num"`
		Type    string `json:"type"`
		Details string `json:"details"`
	}
	if err := json.Unmarshal([]byte(req.Body), &serviceReq); err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid request body"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	if serviceReq.RoomNum <= 0 {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid room number"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	if len(serviceReq.Details) < 5 {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Details must be at least 5 characters"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	var serviceType models.ServiceType
	switch serviceReq.Type {
	case string(models.ServiceTypeCleaning):
		serviceType = models.ServiceTypeCleaning
	case string(models.ServiceTypeFood):
		serviceType = models.ServiceTypeFood
	default:
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid service type"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	ctxWithUserID := context.WithValue(ctx, contextkeys.UserIDKey, userID)
	err := serviceReqSvc.ServiceRequestGetter(ctxWithUserID, serviceReq.RoomNum, serviceType, serviceReq.Details)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusCreated, "Service request submitted", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusCreated,
		Body:       string(body),
	}, nil
}

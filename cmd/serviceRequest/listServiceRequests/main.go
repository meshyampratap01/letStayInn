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
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
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
	lambda.Start(corsmiddleware.WithCORS(handler))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqs, err := serviceReqSvc.GetPendingServiceRequests()
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusInternalServerError, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	dtos := make([]dto.ServiceRequestDTO, 0, len(reqs))
	for _, sr := range reqs {
		dtos = append(dtos, dto.ServiceRequestDTO{
			ID:         sr.ID,
			UserID:     sr.UserID,
			RoomNum:    sr.RoomNum,
			Type:       string(sr.Type),
			Details:    sr.Details,
			Status:     string(sr.Status),
			EmployeeID: sr.AssignedTo,
			IsAssigned: sr.IsAssigned,
		})
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Service requests fetched successfully", dtos)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

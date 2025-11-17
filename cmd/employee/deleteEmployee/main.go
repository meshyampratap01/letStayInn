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
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/managerservice"
)

var managerSvc managerservice.IManagerService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	userRepo := userRepository.NewUserRepo(dynamoDB, tableName)
	managerSvc = managerservice.NewManagerService(userRepo, nil, nil, nil, nil)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	role, ok := roleVal.(string)
	if !ok || role != "Manager" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Manager role required"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
			Body:       string(body),
		}, nil
	}

	employeeID := req.PathParameters["employeeId"]
	if employeeID == "" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Missing employee id"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	err := managerSvc.DeleteEmployeeByID(employeeID)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Employee deleted successfully", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

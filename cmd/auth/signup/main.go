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
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
)

var userSvc userService.IUserService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"

	userRepo := userRepository.NewUserRepo(dynamoDB, tableName)
	userSvc = userService.NewUserService(userRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req dto.SignupRequest

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		resp := response.NewErrorResponse(constants.StatusBadRequest, constants.ErrInvalidRequestBody)
		body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	if err := userSvc.Signup(req, models.RoleGuest); err != nil {
		resp := response.NewErrorResponse(constants.StatusBadRequest, err.Error())
		body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusCreated, "User created successfully", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusCreated,
		Body:       string(body),
	}, nil
}

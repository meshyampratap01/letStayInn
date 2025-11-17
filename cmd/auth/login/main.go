package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
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
	var req dto.LoginRequest

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		resp := response.NewErrorResponse(constants.StatusBadRequest, constants.ErrInvalidRequestBody)
		body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	user, err := userSvc.Login(req.Email, req.Password)
	if err != nil {
		resp := response.NewErrorResponse(constants.StatusUnauthorized, err.Error())
		body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusUnauthorized,
			Body:       string(body),
		}, nil
	}

	token, err := auth.GenerateJWT(user.ID, user.Name, user.Role.String(), constants.JWTExpirationMinute)
	if err != nil {
		resp := response.NewErrorResponse(constants.StatusInternalServerError, constants.ErrFailedToGenerateToken)
		body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Login successful", map[string]string{"token": token})
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

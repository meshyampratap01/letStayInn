package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/db"
	authmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/authMiddleware"
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
	"github.com/meshyampratap01/letStayInn/internal/repository/roomRepository"
	"github.com/meshyampratap01/letStayInn/internal/response"
	"github.com/meshyampratap01/letStayInn/internal/services/roomService"
)

var roomSvc roomService.IRoomService

func init() {
	dynamoDB, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}
	tableName := "letstayinn"
	roomRepo := roomRepository.NewRoomRepo(dynamoDB, tableName)
	roomSvc = roomService.NewRoomService(roomRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authmiddleware.AuthMiddleware(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	roleVal := ctx.Value(contextkeys.UserRoleKey)
	roleStr, ok := roleVal.(string)
	if !ok || roleStr != "Manager" {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusForbidden, "Only managers can update rooms"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusForbidden,
			Body:       string(body),
		}, nil
	}

	roomNumStr := req.PathParameters["roomNum"]
	number, err := strconv.Atoi(roomNumStr)
	if err != nil || number <= 0 {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid room number"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	var reqBody struct {
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		IsAvailable bool    `json:"is_available"`
		Description string  `json:"description"`
	}

	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, "Invalid request body"))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	err = roomSvc.UpdateRoom(number, reqBody.Type, reqBody.Price, reqBody.IsAvailable, reqBody.Description)
	if err != nil {
		body, _ := json.Marshal(response.NewErrorResponse(constants.StatusBadRequest, err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusBadRequest,
			Body:       string(body),
		}, nil
	}

	resp := response.NewSuccessResponse(constants.StatusOK, "Room updated successfully", nil)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

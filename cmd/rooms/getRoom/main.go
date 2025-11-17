package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	authmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/authMiddleware"
	corsmiddleware "github.com/meshyampratap01/letStayInn/internal/middleware/corsMiddleware"
	"github.com/meshyampratap01/letStayInn/internal/models"
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
	if !ok {
		roleStr = ""
	}

	var rooms []models.Room
	var err error

	if roleStr == "Manager" {
		available := req.QueryStringParameters["available"]
		if strings.ToLower(available) == "true" {
			rooms, err = roomSvc.GetAvailableRooms()
		} else {
			rooms, err = roomSvc.GetAllRooms()
		}
	} else {
		rooms, err = roomSvc.GetAvailableRooms()
	}

	if err != nil {
		body, _ := json.Marshal(map[string]string{"message": constants.ErrFailedToFetchUser})
		return events.APIGatewayProxyResponse{
			StatusCode: constants.StatusInternalServerError,
			Body:       string(body),
		}, nil
	}

	roomDTOs := convertModelsToRoomDTOs(rooms)
	resp := response.NewSuccessResponse(constants.StatusOK, "Rooms fetched successfully", roomDTOs)
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: constants.StatusOK,
		Body:       string(body),
	}, nil
}

func convertModelsToRoomDTOs(rooms []models.Room) []dto.RoomDTO {
	roomDTOs := make([]dto.RoomDTO, 0, len(rooms))
	for _, room := range rooms {
		roomDTOs = append(roomDTOs, dto.RoomDTO{
			Number:      room.Number,
			Type:        string(room.Type),
			Price:       room.Price,
			IsAvailable: room.IsAvailable,
			Description: room.Description,
		})
	}
	return roomDTOs
}

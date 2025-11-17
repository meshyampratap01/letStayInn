package roomRepository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type RoomRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewRoomRepo(db *dynamodb.Client, tableName string) IRoomRepository {
	return &RoomRepo{db: db, tableName: tableName}
}

func (r *RoomRepo) GetAllRooms() ([]models.Room, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROOMS"},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Room{}, nil
	}

	var rooms []models.Room
	for _, item := range result.Items {
		var room models.Room
		err := attributevalue.UnmarshalMap(item, &room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepo) GetRoomByNumber(number int) (*models.Room, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND begins_with(sk, ?)"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROOMS"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("room#%d", number)},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("room not found")
	}

	room := &models.Room{}
	err = attributevalue.UnmarshalMap(result.Items[0], room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepo) SaveRoom(room *models.Room) error {
	pk := "ROOMS"
	sk := fmt.Sprintf("room#%d", room.Number)

	statement := "UPDATE " + r.tableName + " SET room_type = ?, price = ?, is_available = ?, description = ? WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: string(room.Type)},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", room.Price)},
			&types.AttributeValueMemberBOOL{Value: room.IsAvailable},
			&types.AttributeValueMemberS{Value: room.Description},
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: sk},
		},
	})

	return err
}

func (r *RoomRepo) SaveRooms(rooms []models.Room) error {
	// for _, room := range rooms {
	// 	if err := r.db.Save(&room).Error; err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (r *RoomRepo) DeleteRoomByNumber(number int) error {
	pk := "ROOMS"
	sk := fmt.Sprintf("room#%d", number)

	statement := "DELETE FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: sk},
		},
	})

	return err
}

func (r *RoomRepo) GetAvailableRooms() ([]models.Room, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND is_available = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROOMS"},
			&types.AttributeValueMemberBOOL{Value: true},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Room{}, nil
	}

	var rooms []models.Room
	for _, item := range result.Items {
		var room models.Room
		err := attributevalue.UnmarshalMap(item, &room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepo) AddRoom(room models.Room) error {
	statement := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'number': ?, 'room_type': ?, 'price': ?, 'is_available': ?, 'description': ?}"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROOMS"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("room#%d", room.Number)},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", room.Number)},
			&types.AttributeValueMemberS{Value: string(room.Type)},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", room.Price)},
			&types.AttributeValueMemberBOOL{Value: room.IsAvailable},
			&types.AttributeValueMemberS{Value: room.Description},
		},
	})

	return err
}

func (r *RoomRepo) GetRoomNumberByBookingID(bookingID string) (string, error) {
	statement := "SELECT room_num FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Bookings"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("Booking#%s", bookingID)},
		},
	})

	if err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("booking not found")
	}

	var booking models.Booking
	err = attributevalue.UnmarshalMap(result.Items[0], &booking)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", booking.RoomNum), nil
}

func (r *RoomRepo) RoomExists(number int) (bool, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND begins_with(sk, ?)"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROOMS"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("room#%d", number)},
		},
	})

	if err != nil {
		return false, err
	}

	return len(result.Items) > 0, nil
}

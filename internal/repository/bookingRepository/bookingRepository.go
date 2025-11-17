package bookingRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type BookingRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewBookingRepo(db *dynamodb.Client, tableName string) BookingRepository {
	return &BookingRepo{db: db, tableName: tableName}
}

func (r *BookingRepo) GetAllBookings() ([]models.Booking, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Bookings"},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Booking{}, nil
	}

	var bookings []models.Booking
	for _, item := range result.Items {
		type tempBooking struct {
			ID       string `dynamodbav:"id"`
			UserID   string `dynamodbav:"user_id"`
			RoomID   string `dynamodbav:"room_id"`
			RoomNum  int    `dynamodbav:"room_num"`
			CheckIn  string `dynamodbav:"check_in"`
			CheckOut string `dynamodbav:"check_out"`
			Status   string `dynamodbav:"status"`
			FoodReq  bool   `dynamodbav:"food_req"`
			CleanReq bool   `dynamodbav:"clean_req"`
		}

		var temp tempBooking
		err := attributevalue.UnmarshalMap(item, &temp)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05.999999999 -0700 MST"
		checkInTime, _ := time.Parse(layout, temp.CheckIn)
		checkOutTime, _ := time.Parse(layout, temp.CheckOut)

		booking := models.Booking{
			ID:       temp.ID,
			UserID:   temp.UserID,
			RoomID:   temp.RoomID,
			RoomNum:  temp.RoomNum,
			CheckIn:  checkInTime,
			CheckOut: checkOutTime,
			Status:   temp.Status,
			FoodReq:  temp.FoodReq,
			CleanReq: temp.CleanReq,
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepo) SaveBooking(booking models.Booking) error {
	userPK := fmt.Sprintf("User#%s", booking.UserID)
	userSK := fmt.Sprintf("booking#%s", booking.ID)

	statement1 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'user_id': ?, 'room_id': ?, 'room_num': ?, 'check_in': ?, 'check_out': ?, 'status': ?, 'food_req': ?, 'clean_req': ?}"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: userPK},
			&types.AttributeValueMemberS{Value: userSK},
			&types.AttributeValueMemberS{Value: booking.ID},
			&types.AttributeValueMemberS{Value: booking.UserID},
			&types.AttributeValueMemberS{Value: booking.RoomID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", booking.RoomNum)},
			&types.AttributeValueMemberS{Value: booking.CheckIn.String()},
			&types.AttributeValueMemberS{Value: booking.CheckOut.String()},
			&types.AttributeValueMemberS{Value: booking.Status},
			&types.AttributeValueMemberBOOL{Value: booking.FoodReq},
			&types.AttributeValueMemberBOOL{Value: booking.CleanReq},
		},
	})

	if err != nil {
		return err
	}

	bookingsPK := "Bookings"
	bookingsSK := fmt.Sprintf("booking#%s", booking.ID)

	statement2 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'user_id': ?, 'room_id': ?, 'room_num': ?, 'check_in': ?, 'check_out': ?, 'status': ?, 'food_req': ?, 'clean_req': ?}"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: bookingsPK},
			&types.AttributeValueMemberS{Value: bookingsSK},
			&types.AttributeValueMemberS{Value: booking.ID},
			&types.AttributeValueMemberS{Value: booking.UserID},
			&types.AttributeValueMemberS{Value: booking.RoomID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", booking.RoomNum)},
			&types.AttributeValueMemberS{Value: booking.CheckIn.String()},
			&types.AttributeValueMemberS{Value: booking.CheckOut.String()},
			&types.AttributeValueMemberS{Value: booking.Status},
			&types.AttributeValueMemberBOOL{Value: booking.FoodReq},
			&types.AttributeValueMemberBOOL{Value: booking.CleanReq},
		},
	})

	return err
}

func (r *BookingRepo) GetBookingsByUserID(userID string) ([]models.Booking, error) {
	pk := fmt.Sprintf("User#%s", userID)
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Booking{}, nil
	}

	var bookings []models.Booking
	for _, item := range result.Items {
		// Use a temporary struct to unmarshal with DynamoDB types
		type tempBooking struct {
			ID       string `dynamodbav:"id"`
			UserID   string `dynamodbav:"user_id"`
			RoomID   string `dynamodbav:"room_id"`
			RoomNum  int    `dynamodbav:"room_num"`
			CheckIn  string `dynamodbav:"check_in"`
			CheckOut string `dynamodbav:"check_out"`
			Status   string `dynamodbav:"status"`
			FoodReq  bool   `dynamodbav:"food_req"`
			CleanReq bool   `dynamodbav:"clean_req"`
		}

		var temp tempBooking
		err := attributevalue.UnmarshalMap(item, &temp)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05.999999999 -0700 MST"
		checkInTime, _ := time.Parse(layout, temp.CheckIn)
		checkOutTime, _ := time.Parse(layout, temp.CheckOut)

		booking := models.Booking{
			ID:       temp.ID,
			UserID:   temp.UserID,
			RoomID:   temp.RoomID,
			RoomNum:  temp.RoomNum,
			CheckIn:  checkInTime,
			CheckOut: checkOutTime,
			Status:   temp.Status,
			FoodReq:  temp.FoodReq,
			CleanReq: temp.CleanReq,
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepo) UpdateBooking(updated models.Booking) error {
	userPK := fmt.Sprintf("User#%s", updated.UserID)
	userSK := fmt.Sprintf("booking#%s", updated.ID)

	statement1 := "UPDATE " + r.tableName + " SET status = ?, food_req = ?, clean_req = ?, check_in = ?, check_out = ? WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: updated.Status},
			&types.AttributeValueMemberBOOL{Value: updated.FoodReq},
			&types.AttributeValueMemberBOOL{Value: updated.CleanReq},
			&types.AttributeValueMemberS{Value: updated.CheckIn.String()},
			&types.AttributeValueMemberS{Value: updated.CheckOut.String()},
			&types.AttributeValueMemberS{Value: userPK},
			&types.AttributeValueMemberS{Value: userSK},
		},
	})

	if err != nil {
		return err
	}

	bookingsSK := fmt.Sprintf("booking#%s", updated.ID)

	statement2 := "UPDATE " + r.tableName + " SET status = ?, food_req = ?, clean_req = ?, check_in = ?, check_out = ? WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: updated.Status},
			&types.AttributeValueMemberBOOL{Value: updated.FoodReq},
			&types.AttributeValueMemberBOOL{Value: updated.CleanReq},
			&types.AttributeValueMemberS{Value: updated.CheckIn.String()},
			&types.AttributeValueMemberS{Value: updated.CheckOut.String()},
			&types.AttributeValueMemberS{Value: "Bookings"},
			&types.AttributeValueMemberS{Value: bookingsSK},
		},
	})

	return err
}

func (r *BookingRepo) GetBookingByID(bookingID string) (*models.Booking, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Bookings"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("booking#%s", bookingID)},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	type tempBooking struct {
		ID       string `dynamodbav:"id"`
		UserID   string `dynamodbav:"user_id"`
		RoomID   string `dynamodbav:"room_id"`
		RoomNum  int    `dynamodbav:"room_num"`
		CheckIn  string `dynamodbav:"check_in"`
		CheckOut string `dynamodbav:"check_out"`
		Status   string `dynamodbav:"status"`
		FoodReq  bool   `dynamodbav:"food_req"`
		CleanReq bool   `dynamodbav:"clean_req"`
	}

	var temp tempBooking
	err = attributevalue.UnmarshalMap(result.Items[0], &temp)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	checkInTime, _ := time.Parse(layout, temp.CheckIn)
	checkOutTime, _ := time.Parse(layout, temp.CheckOut)

	booking := &models.Booking{
		ID:       temp.ID,
		UserID:   temp.UserID,
		RoomID:   temp.RoomID,
		RoomNum:  temp.RoomNum,
		CheckIn:  checkInTime,
		CheckOut: checkOutTime,
		Status:   temp.Status,
		FoodReq:  temp.FoodReq,
		CleanReq: temp.CleanReq,
	}

	return booking, nil
}

func (r *BookingRepo) GetActiveBookings() ([]models.Booking, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND #s <> ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Bookings"},
			&types.AttributeValueMemberS{Value: models.BookingStatusCancelled},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Booking{}, nil
	}

	var bookings []models.Booking
	for _, item := range result.Items {
		type tempBooking struct {
			ID       string `dynamodbav:"id"`
			UserID   string `dynamodbav:"user_id"`
			RoomID   string `dynamodbav:"room_id"`
			RoomNum  int    `dynamodbav:"room_num"`
			CheckIn  string `dynamodbav:"check_in"`
			CheckOut string `dynamodbav:"check_out"`
			Status   string `dynamodbav:"status"`
			FoodReq  bool   `dynamodbav:"food_req"`
			CleanReq bool   `dynamodbav:"clean_req"`
		}

		var temp tempBooking
		err := attributevalue.UnmarshalMap(item, &temp)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05.999999999 -0700 MST"
		checkInTime, _ := time.Parse(layout, temp.CheckIn)
		checkOutTime, _ := time.Parse(layout, temp.CheckOut)

		booking := models.Booking{
			ID:       temp.ID,
			UserID:   temp.UserID,
			RoomID:   temp.RoomID,
			RoomNum:  temp.RoomNum,
			CheckIn:  checkInTime,
			CheckOut: checkOutTime,
			Status:   temp.Status,
			FoodReq:  temp.FoodReq,
			CleanReq: temp.CleanReq,
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepo) CheckRoomBooked(roomNumber int) (bool, error) {
	statement := "SELECT COUNT(*) as cnt FROM " + r.tableName + " WHERE pk = ? AND room_num = ? AND #s = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Bookings"},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", roomNumber)},
			&types.AttributeValueMemberS{Value: models.BookingStatusBooked},
		},
	})

	if err != nil {
		return false, err
	}

	if len(result.Items) == 0 {
		return false, nil
	}

	type countResult struct {
		Count int `dynamodbav:"cnt"`
	}

	var cr countResult
	err = attributevalue.UnmarshalMap(result.Items[0], &cr)
	if err != nil {
		return false, err
	}

	return cr.Count > 0, nil
}

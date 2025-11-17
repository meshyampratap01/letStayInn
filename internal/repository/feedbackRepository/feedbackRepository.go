package feedbackRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type feedbackRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewFeedbackRepo(db *dynamodb.Client, tableName string) FeedbackRepository {
	return &feedbackRepo{db: db, tableName: tableName}
}

func (r *feedbackRepo) SaveFeedback(f models.Feedback) error {
	pk := "Feedbacks"
	sk := fmt.Sprintf("feedback#%s", f.ID)

	statement := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'user_id': ?, 'user_name': ?, 'message': ?, 'booking_id': ?, 'room_num': ?, 'rating': ?, 'created_at': ?}"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: sk},
			&types.AttributeValueMemberS{Value: f.ID},
			&types.AttributeValueMemberS{Value: f.UserID},
			&types.AttributeValueMemberS{Value: f.UserName},
			&types.AttributeValueMemberS{Value: f.Message},
			&types.AttributeValueMemberS{Value: f.BookingID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", f.RoomNum)},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", f.Rating)},
			&types.AttributeValueMemberS{Value: f.CreatedAt.String()},
		},
	})

	return err
}

func (r *feedbackRepo) GetAllFeedback() ([]models.Feedback, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Feedbacks"},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.Feedback{}, nil
	}

	var feedbacks []models.Feedback
	for _, item := range result.Items {
		type tempFeedback struct {
			ID        string `dynamodbav:"id"`
			UserID    string `dynamodbav:"user_id"`
			UserName  string `dynamodbav:"user_name"`
			Message   string `dynamodbav:"message"`
			BookingID string `dynamodbav:"booking_id"`
			RoomNum   int    `dynamodbav:"room_num"`
			Rating    int    `dynamodbav:"rating"`
			CreatedAt string `dynamodbav:"created_at"`
		}

		var temp tempFeedback
		err := attributevalue.UnmarshalMap(item, &temp)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05.999999999 -0700 MST"
		createdAt, _ := time.Parse(layout, temp.CreatedAt)

		feedback := models.Feedback{
			ID:        temp.ID,
			UserID:    temp.UserID,
			UserName:  temp.UserName,
			Message:   temp.Message,
			BookingID: temp.BookingID,
			RoomNum:   temp.RoomNum,
			Rating:    temp.Rating,
			CreatedAt: createdAt,
		}
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

func (r *feedbackRepo) DeleteFeedback(feedbackID string) error {
	statement := "DELETE FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Feedbacks"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("feedback#%s", feedbackID)},
		},
	})

	return err
}

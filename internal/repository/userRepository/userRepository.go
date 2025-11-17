package userRepository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type UserRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewUserRepo(db *dynamodb.Client, tableName string) *UserRepo {
	return &UserRepo{db: db, tableName: tableName}
}

func (r *UserRepo) GetAllUsers() ([]models.User, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Users"},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.User{}, nil
	}

	var users []models.User
	for _, item := range result.Items {
		var user models.User
		err := attributevalue.UnmarshalMap(item, &user)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) FindUserByEmail(users []models.User, email string) *models.User {
	for _, u := range users {
		if u.Email == email {
			return &u
		}
	}
	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: email},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf(constants.ErrUserNotFound)
	}

	user := &models.User{}
	err = attributevalue.UnmarshalMap(result.Items[0], user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) SaveUser(newUser models.User) error {
	statement1 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'name': ?, 'email': ?, 'password': ?, 'role': ?, 'available': ?}"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: newUser.Email},
			&types.AttributeValueMemberS{Value: newUser.ID},
			&types.AttributeValueMemberS{Value: newUser.ID},
			&types.AttributeValueMemberS{Value: newUser.Name},
			&types.AttributeValueMemberS{Value: newUser.Email},
			&types.AttributeValueMemberS{Value: newUser.Password},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", newUser.Role)},
			&types.AttributeValueMemberBOOL{Value: newUser.Available},
		},
	})

	if err != nil {
		return err
	}

	userSK := fmt.Sprintf("user#%s", newUser.ID)

	statement2 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'name': ?, 'email': ?, 'password': ?, 'role': ?, 'available': ?}"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Users"},
			&types.AttributeValueMemberS{Value: userSK},
			&types.AttributeValueMemberS{Value: newUser.ID},
			&types.AttributeValueMemberS{Value: newUser.Name},
			&types.AttributeValueMemberS{Value: newUser.Email},
			&types.AttributeValueMemberS{Value: newUser.Password},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", newUser.Role)},
			&types.AttributeValueMemberBOOL{Value: newUser.Available},
		},
	})

	return err
}

func (r *UserRepo) UpdateUser(user *models.User) error {
	statement1 := "UPDATE " + r.tableName + " SET name = ?, password = ?, role = ?, available = ? WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: user.Name},
			&types.AttributeValueMemberS{Value: user.Password},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.Role)},
			&types.AttributeValueMemberBOOL{Value: user.Available},
			&types.AttributeValueMemberS{Value: user.Email},
			&types.AttributeValueMemberS{Value: user.ID},
		},
	})

	if err != nil {
		return err
	}

	userSK := fmt.Sprintf("user#%s", user.ID)
	statement2 := "UPDATE " + r.tableName + " SET name = ?, password = ?, role = ?, available = ? WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: user.Name},
			&types.AttributeValueMemberS{Value: user.Password},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.Role)},
			&types.AttributeValueMemberBOOL{Value: user.Available},
			&types.AttributeValueMemberS{Value: "Users"},
			&types.AttributeValueMemberS{Value: userSK},
		},
	})

	return err
}

func (r *UserRepo) SaveAllUsers(users []models.User) error {
	// for _, user := range users {
	// 	if err := r.db.Save(&user).Error; err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (r *UserRepo) GetUserByID(userID string) (*models.User, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Users"},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("user#%s", userID)},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf(constants.ErrUserNotFound)
	}

	user := &models.User{}
	err = attributevalue.UnmarshalMap(result.Items[0], user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) ToggleStaffAvailability(userID string) error {
	var user models.User
	// if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
	// 	return err
	// }
	user.Available = !user.Available
	// return r.db.Save(&user).Error
	return nil
}

func (r *UserRepo) GetStaffAvailability(userID string) (bool, error) {
	var user models.User
	// if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
	// 	return false, err
	// }
	return user.Available, nil
}

func (r *UserRepo) DeleteUserByID(userID string) error {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}

	statement1 := "DELETE FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: user.Email},
			&types.AttributeValueMemberS{Value: userID},
		},
	})

	if err != nil {
		return err
	}

	userSK := fmt.Sprintf("user#%s", userID)
	statement2 := "DELETE FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "Users"},
			&types.AttributeValueMemberS{Value: userSK},
		},
	})

	return err
}

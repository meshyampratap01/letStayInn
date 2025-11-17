package serviceRequestRepository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type ServiceRequestRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewServiceRequestRepo(db *dynamodb.Client, tableName string) ServiceRequestRepository {
	return &ServiceRequestRepo{db: db, tableName: tableName}
}

func (r *ServiceRequestRepo) SaveServiceRequest(req models.ServiceRequest) error {
	// Partition 1: ServiceRequests - PK=ServiceRequests, SK=service#{id}
	statement1 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'userID': ?, 'bookingID': ?, 'roomNum': ?, 'type': ?, 'status': ?, 'isAssigned': ?, 'assignedTo': ?, 'details': ?}"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: "service#" + req.ID},
			&types.AttributeValueMemberS{Value: req.ID},
			&types.AttributeValueMemberS{Value: req.UserID},
			&types.AttributeValueMemberS{Value: req.BookingID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: string(req.Type)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberBOOL{Value: req.IsAssigned},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Details},
		},
	})

	if err != nil {
		return err
	}

	// Partition 2: Room - PK={roomNum}, SK=service#{id}
	statement2 := "INSERT INTO " + r.tableName + " VALUE {'pk': ?, 'sk': ?, 'id': ?, 'userID': ?, 'bookingID': ?, 'roomNum': ?, 'type': ?, 'status': ?, 'isAssigned': ?, 'assignedTo': ?, 'details': ?}"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: "service#" + req.ID},
			&types.AttributeValueMemberS{Value: req.ID},
			&types.AttributeValueMemberS{Value: req.UserID},
			&types.AttributeValueMemberS{Value: req.BookingID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: string(req.Type)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberBOOL{Value: req.IsAssigned},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Details},
		},
	})

	return err
}

func (r *ServiceRequestRepo) LoadServiceRequests() ([]models.ServiceRequest, error) {
	// Query ServiceRequests partition - PK=ServiceRequests
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.ServiceRequest{}, nil
	}

	var requests []models.ServiceRequest

	for _, item := range result.Items {
		var tempReq struct {
			PK         string `dynamodbav:"pk"`
			SK         string `dynamodbav:"sk"`
			ID         string `dynamodbav:"id"`
			UserID     string `dynamodbav:"userID"`
			BookingID  string `dynamodbav:"bookingID"`
			RoomNum    int    `dynamodbav:"roomNum"`
			Type       string `dynamodbav:"type"`
			Status     string `dynamodbav:"status"`
			IsAssigned bool   `dynamodbav:"isAssigned"`
			AssignedTo string `dynamodbav:"assignedTo"`
			Details    string `dynamodbav:"details"`
		}

		err := attributevalue.UnmarshalMap(item, &tempReq)
		if err != nil {
			continue
		}

		req := models.ServiceRequest{
			PK:         tempReq.PK,
			SK:         tempReq.SK,
			ID:         tempReq.ID,
			UserID:     tempReq.UserID,
			BookingID:  tempReq.BookingID,
			RoomNum:    tempReq.RoomNum,
			Type:       models.ServiceType(tempReq.Type),
			Status:     models.ServiceStatus(tempReq.Status),
			IsAssigned: tempReq.IsAssigned,
			AssignedTo: tempReq.AssignedTo,
			Details:    tempReq.Details,
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *ServiceRequestRepo) GetPendingServiceRequests() ([]models.ServiceRequest, error) {
	// Query ServiceRequests partition for pending items - PK=ServiceRequests, status=Pending
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND status = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: string(models.ServiceStatusPending)},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.ServiceRequest{}, nil
	}

	var requests []models.ServiceRequest

	for _, item := range result.Items {
		var tempReq struct {
			PK         string `dynamodbav:"pk"`
			SK         string `dynamodbav:"sk"`
			ID         string `dynamodbav:"id"`
			UserID     string `dynamodbav:"userID"`
			BookingID  string `dynamodbav:"bookingID"`
			RoomNum    int    `dynamodbav:"roomNum"`
			Type       string `dynamodbav:"type"`
			Status     string `dynamodbav:"status"`
			IsAssigned bool   `dynamodbav:"isAssigned"`
			AssignedTo string `dynamodbav:"assignedTo"`
			Details    string `dynamodbav:"details"`
		}

		err := attributevalue.UnmarshalMap(item, &tempReq)
		if err != nil {
			continue
		}

		req := models.ServiceRequest{
			PK:         tempReq.PK,
			SK:         tempReq.SK,
			ID:         tempReq.ID,
			UserID:     tempReq.UserID,
			BookingID:  tempReq.BookingID,
			RoomNum:    tempReq.RoomNum,
			Type:       models.ServiceType(tempReq.Type),
			Status:     models.ServiceStatus(tempReq.Status),
			IsAssigned: tempReq.IsAssigned,
			AssignedTo: tempReq.AssignedTo,
			Details:    tempReq.Details,
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *ServiceRequestRepo) GetServiceRequestByReqID(id string) (*models.ServiceRequest, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND sk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: "service#" + id},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("service request not found")
	}

	var tempReq struct {
		PK         string `dynamodbav:"pk"`
		SK         string `dynamodbav:"sk"`
		ID         string `dynamodbav:"id"`
		UserID     string `dynamodbav:"userID"`
		BookingID  string `dynamodbav:"bookingID"`
		RoomNum    int    `dynamodbav:"roomNum"`
		Type       string `dynamodbav:"type"`
		Status     string `dynamodbav:"status"`
		IsAssigned bool   `dynamodbav:"isAssigned"`
		AssignedTo string `dynamodbav:"assignedTo"`
		Details    string `dynamodbav:"details"`
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &tempReq)
	if err != nil {
		return nil, err
	}

	req := &models.ServiceRequest{
		PK:         tempReq.PK,
		SK:         tempReq.SK,
		ID:         tempReq.ID,
		UserID:     tempReq.UserID,
		BookingID:  tempReq.BookingID,
		RoomNum:    tempReq.RoomNum,
		Type:       models.ServiceType(tempReq.Type),
		Status:     models.ServiceStatus(tempReq.Status),
		IsAssigned: tempReq.IsAssigned,
		AssignedTo: tempReq.AssignedTo,
		Details:    tempReq.Details,
	}

	return req, nil
}

func (r *ServiceRequestRepo) GetServiceRequestByRoomNum(roomNum int) (*models.ServiceRequest, error) {
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: fmt.Sprintf("%d", roomNum)},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no service request found for room %d", roomNum)
	}

	var tempReq struct {
		PK         string `dynamodbav:"pk"`
		SK         string `dynamodbav:"sk"`
		ID         string `dynamodbav:"id"`
		UserID     string `dynamodbav:"userID"`
		BookingID  string `dynamodbav:"bookingID"`
		RoomNum    int    `dynamodbav:"roomNum"`
		Type       string `dynamodbav:"type"`
		Status     string `dynamodbav:"status"`
		IsAssigned bool   `dynamodbav:"isAssigned"`
		AssignedTo string `dynamodbav:"assignedTo"`
		Details    string `dynamodbav:"details"`
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &tempReq)
	if err != nil {
		return nil, err
	}

	req := &models.ServiceRequest{
		PK:         tempReq.PK,
		SK:         tempReq.SK,
		ID:         tempReq.ID,
		UserID:     tempReq.UserID,
		BookingID:  tempReq.BookingID,
		RoomNum:    tempReq.RoomNum,
		Type:       models.ServiceType(tempReq.Type),
		Status:     models.ServiceStatus(tempReq.Status),
		IsAssigned: tempReq.IsAssigned,
		AssignedTo: tempReq.AssignedTo,
		Details:    tempReq.Details,
	}

	return req, nil
}

func (r *ServiceRequestRepo) UpdateServiceRequest(req *models.ServiceRequest) error {
	// Update ServiceRequests partition
	statement1 := "UPDATE " + r.tableName + " SET userID = ?, bookingID = ?, roomNum = ?, type = ?, status = ?, isAssigned = ?, assignedTo = ?, details = ? WHERE pk = ? AND sk = ?"

	_, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: req.UserID},
			&types.AttributeValueMemberS{Value: req.BookingID},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: string(req.Type)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberBOOL{Value: req.IsAssigned},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Details},
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: "service#" + req.ID},
		},
	})

	if err != nil {
		return err
	}

	// Update Room partition
	statement2 := "UPDATE " + r.tableName + " SET userID = ?, bookingID = ?, type = ?, status = ?, isAssigned = ?, assignedTo = ?, details = ? WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: req.UserID},
			&types.AttributeValueMemberS{Value: req.BookingID},
			&types.AttributeValueMemberS{Value: string(req.Type)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberBOOL{Value: req.IsAssigned},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Details},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: "service#" + req.ID},
		},
	})

	return err
}

func (r *ServiceRequestRepo) GetAssignedServiceRequests(employeeID string) ([]models.ServiceRequest, error) {
	// Query ServiceRequests partition for items assigned to specific employee
	statement := "SELECT * FROM " + r.tableName + " WHERE pk = ? AND assignedTo = ?"

	result, err := r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: employeeID},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return []models.ServiceRequest{}, nil
	}

	var requests []models.ServiceRequest

	for _, item := range result.Items {
		var tempReq struct {
			PK         string `dynamodbav:"pk"`
			SK         string `dynamodbav:"sk"`
			ID         string `dynamodbav:"id"`
			UserID     string `dynamodbav:"userID"`
			BookingID  string `dynamodbav:"bookingID"`
			RoomNum    int    `dynamodbav:"roomNum"`
			Type       string `dynamodbav:"type"`
			Status     string `dynamodbav:"status"`
			IsAssigned bool   `dynamodbav:"isAssigned"`
			AssignedTo string `dynamodbav:"assignedTo"`
			Details    string `dynamodbav:"details"`
		}

		err := attributevalue.UnmarshalMap(item, &tempReq)
		if err != nil {
			continue
		}

		req := models.ServiceRequest{
			PK:         tempReq.PK,
			SK:         tempReq.SK,
			ID:         tempReq.ID,
			UserID:     tempReq.UserID,
			BookingID:  tempReq.BookingID,
			RoomNum:    tempReq.RoomNum,
			Type:       models.ServiceType(tempReq.Type),
			Status:     models.ServiceStatus(tempReq.Status),
			IsAssigned: tempReq.IsAssigned,
			AssignedTo: tempReq.AssignedTo,
			Details:    tempReq.Details,
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *ServiceRequestRepo) UpdateIsAssigned(reqID string, employeeID string, isAssigned bool) error {
	// Get current request first
	req, err := r.GetServiceRequestByReqID(reqID)
	if err != nil {
		return err
	}

	// Update isAssigned field in ServiceRequests partition
	statement1 := "UPDATE " + r.tableName + " SET isAssigned = ?, assignedTo = ? WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement1,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberBOOL{Value: isAssigned},
			&types.AttributeValueMemberS{Value: employeeID},
			&types.AttributeValueMemberS{Value: "ServiceRequests"},
			&types.AttributeValueMemberS{Value: "service#" + reqID},
		},
	})

	if err != nil {
		return err
	}

	// Update Room partition
	statement2 := "UPDATE " + r.tableName + " SET isAssigned = ?, assignedTo = ? WHERE pk = ? AND sk = ?"

	_, err = r.db.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement2,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberBOOL{Value: isAssigned},
			&types.AttributeValueMemberS{Value: employeeID},
			&types.AttributeValueMemberS{Value: fmt.Sprintf("%d", req.RoomNum)},
			&types.AttributeValueMemberS{Value: "service#" + reqID},
		},
	})

	return err
}

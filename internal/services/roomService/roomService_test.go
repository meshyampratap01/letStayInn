package roomService

// import (
// 	"errors"
// 	"testing"

// 	"github.com/meshyampratap01/letStayInn/internal/mocks"
// 	"github.com/meshyampratap01/letStayInn/internal/models"
// 	"go.uber.org/mock/gomock"
// )

// func setupRoomService(ctrl *gomock.Controller) (*RoomService, *mocks.MockIRoomRepository) {
// 	mockRoomRepo := mocks.NewMockIRoomRepository(ctrl)
// 	service := NewRoomService(mockRoomRepo)
// 	return service.(*RoomService), mockRoomRepo
// }

// func TestGetTotalRooms_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{}, {}}, nil)
// 	total, err := service.GetTotalRooms()
// 	if err != nil || total != 2 {
// 		t.Errorf("expected 2 rooms, got %v, err %v", total, err)
// 	}
// }

// func TestGetTotalRooms_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(nil, errors.New("db error"))
// 	_, err := service.GetTotalRooms()
// 	if err == nil || err.Error() != "db error" {
// 		t.Errorf("expected db error, got %v", err)
// 	}
// }

// func TestGetAvailableRooms_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{IsAvailable: true}}, nil)
// 	rooms, err := service.GetAvailableRooms()
// 	if err != nil || len(rooms) != 1 {
// 		t.Errorf("expected 1 available room, got %v, err %v", len(rooms), err)
// 	}
// }

// func TestGetAvailableRooms_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAvailableRooms().Return(nil, errors.New("db error"))
// 	_, err := service.GetAvailableRooms()
// 	if err == nil || err.Error() != "db error" {
// 		t.Errorf("expected db error, got %v", err)
// 	}
// }

// func TestGetTotalAvailableRooms_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAvailableRooms().Return([]models.Room{{}, {}}, nil)
// 	total, err := service.GetTotalAvailableRooms()
// 	if err != nil || total != 2 {
// 		t.Errorf("expected 2 available rooms, got %v, err %v", total, err)
// 	}
// }

// func TestGetTotalAvailableRooms_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAvailableRooms().Return(nil, errors.New("db error"))
// 	_, err := service.GetTotalAvailableRooms()
// 	if err == nil || err.Error() != "db error" {
// 		t.Errorf("expected db error, got %v", err)
// 	}
// }

// func TestGetAllRooms_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{Number: 101}}, nil)
// 	rooms, err := service.GetAllRooms()
// 	if err != nil || len(rooms) != 1 {
// 		t.Errorf("expected 1 room, got %v, err %v", len(rooms), err)
// 	}
// }

// func TestGetAllRooms_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(nil, errors.New("db error"))
// 	_, err := service.GetAllRooms()
// 	if err == nil || err.Error() != "db error" {
// 		t.Errorf("expected db error, got %v", err)
// 	}
// }

// func TestAddRoom_ValidationErrors(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, _ := setupRoomService(ctrl)
// 	if err := service.AddRoom(0, "Deluxe", 100, true, "desc"); err == nil {
// 		t.Errorf("expected invalid room number error")
// 	}
// 	if err := service.AddRoom(101, "", 100, true, "desc"); err == nil {
// 		t.Errorf("expected invalid room type error")
// 	}
// 	if err := service.AddRoom(101, "Deluxe", 0, true, "desc"); err == nil {
// 		t.Errorf("expected invalid room price error")
// 	}
// 	if err := service.AddRoom(101, "Deluxe", 100, true, ""); err == nil {
// 		t.Errorf("expected empty room desc error")
// 	}
// }

// func TestAddRoom_DuplicateRoomNumber(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().RoomExists(101).Return(true, nil)
// 	if err := service.AddRoom(101, "Deluxe", 100, true, "desc"); err == nil || err.Error() != "room number already exists" {
// 		t.Errorf("expected duplicate room number error, got %v", err)
// 	}
// }

// func TestAddRoom_RepoError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().RoomExists(101).Return(false, errors.New("repo error"))
// 	if err := service.AddRoom(101, "Deluxe", 100, true, "desc"); err == nil || err.Error() != "repo error" {
// 		t.Errorf("expected repo error, got %v", err)
// 	}
// }

// func TestAddRoom_AddRoomError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().RoomExists(101).Return(false, nil)
// 	mockRoomRepo.EXPECT().AddRoom(gomock.Any()).Return(errors.New("add error"))
// 	if err := service.AddRoom(101, "Deluxe", 100, true, "desc"); err == nil || err.Error() != "add error" {
// 		t.Errorf("expected add error, got %v", err)
// 	}
// }

// func TestUpdateRoom_ChoiceValidation(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, _ := setupRoomService(ctrl)
// 	if err := service.UpdateRoom(101, 0, "Deluxe", 100, true, "desc"); err == nil {
// 		t.Errorf("expected invalid choice error")
// 	}
// 	if err := service.UpdateRoom(101, 5, "Deluxe", 100, true, "desc"); err == nil {
// 		t.Errorf("expected invalid choice error")
// 	}
// }

// func TestUpdateRoom_RoomNotFound(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{}, nil)
// 	if err := service.UpdateRoom(101, 1, "Deluxe", 100, true, "desc"); err == nil || err.Error() != "room not found: room number 101" {
// 		t.Errorf("expected room not found error, got %v", err)
// 	}
// }

// func TestUpdateRoom_PersistenceError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().GetAllRooms().Return([]models.Room{{Number: 101}}, nil)
// 	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(errors.New("save error"))
// 	if err := service.UpdateRoom(101, 1, "Deluxe", 100, true, "desc"); err == nil || err.Error() != "failed to save rooms: save error" {
// 		t.Errorf("expected save error, got %v", err)
// 	}
// }

// func TestDeleteRoom_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().DeleteRoomByNumber(101).Return(nil)
// 	if err := service.DeleteRoom(101); err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestDeleteRoom_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().DeleteRoomByNumber(101).Return(errors.New("delete error"))
// 	if err := service.DeleteRoom(101); err == nil || err.Error() != "delete error" {
// 		t.Errorf("expected delete error, got %v", err)
// 	}
// }

// func TestRoomExists_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().RoomExists(101).Return(true, nil)
// 	exists, err := service.RoomExists(101)
// 	if err != nil || !exists {
// 		t.Errorf("expected room to exist, got %v, err %v", exists, err)
// 	}
// }

// func TestRoomExists_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	mockRoomRepo.EXPECT().RoomExists(101).Return(false, errors.New("exists error"))
// 	_, err := service.RoomExists(101)
// 	if err == nil || err.Error() != "exists error" {
// 		t.Errorf("expected exists error, got %v", err)
// 	}
// }

// func TestUpdateRoom_TypeSuccess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Type: "Standard"}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil)
// 	err := service.UpdateRoom(101, 1, "Deluxe", 100, true, "desc")
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestUpdateRoom_TypeEmptyError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Type: "Standard"}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	err := service.UpdateRoom(101, 1, "", 100, true, "desc")
// 	if err == nil || err.Error() != "invalid input: room type cannot be empty" {
// 		t.Errorf("expected room type cannot be empty error, got %v", err)
// 	}
// }

// func TestUpdateRoom_PriceSuccess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Price: 50}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil)
// 	err := service.UpdateRoom(101, 2, "Deluxe", 100, true, "desc")
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestUpdateRoom_PriceZeroError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Price: 50}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	err := service.UpdateRoom(101, 2, "Deluxe", 0, true, "desc")
// 	if err == nil || err.Error() != "invalid input: price 0.00 must be greater than 0" {
// 		t.Errorf("expected price must be greater than 0 error, got %v", err)
// 	}
// }

// func TestUpdateRoom_AvailabilitySuccess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, IsAvailable: false}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil)
// 	err := service.UpdateRoom(101, 3, "Deluxe", 100, true, "desc")
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestUpdateRoom_DescriptionSuccess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Description: "old"}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	mockRoomRepo.EXPECT().SaveRooms(gomock.Any()).Return(nil)
// 	err := service.UpdateRoom(101, 4, "Deluxe", 100, true, "new description")
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestUpdateRoom_DescriptionShortError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	service, mockRoomRepo := setupRoomService(ctrl)
// 	rooms := []models.Room{{Number: 101, Description: "old"}}
// 	mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)
// 	err := service.UpdateRoom(101, 4, "Deluxe", 100, true, "abc")
// 	if err == nil || err.Error() != "invalid input: description must be at least 5 characters long" {
// 		t.Errorf("expected description length error, got %v", err)
// 	}
// }

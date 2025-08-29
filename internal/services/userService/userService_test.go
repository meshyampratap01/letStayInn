package userService

import (
	"errors"
	"testing"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"go.uber.org/mock/gomock"
)

func setupUserService(ctrl *gomock.Controller) (*UserService, *mocks.MockUserRepository) {
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockUserRepo)
	return service.(*UserService), mockUserRepo
}

func TestSignup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
	msg, err := service.Signup("John", "john@example.com", "pass", int(models.RoleGuest))
	if err != nil || msg == "" {
		t.Errorf("expected success, got %v, msg %v", err, msg)
	}
}

func TestSignup_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := setupUserService(ctrl)
	_, err := service.Signup("John", "john@example.com", "pass", -1)
	if err == nil || err.Error() != "invalid role" {
		t.Errorf("expected invalid role error, got %v", err)
	}
}

func TestGetAllEmployees_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	users := []models.User{{Role: models.RoleKitchenStaff}, {Role: models.RoleManager}, {Role: models.RoleGuest}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	emps, err := service.GetAllEmployees()
	if err != nil || len(emps) != 2 {
		t.Errorf("expected 2 employees, got %v, err %v", len(emps), err)
	}
}

func TestGetAllEmployees_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	_, err := service.GetAllEmployees()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{{}}, nil)
	users, err := service.GetAllUsers()
	if err != nil || len(users) != 1 {
		t.Errorf("expected 1 user, got %v, err %v", len(users), err)
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
	_, err := service.GetAllUsers()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	user := &models.User{ID: "u1"}
	mockUserRepo.EXPECT().GetUserByID("u1").Return(user, nil)
	u, err := service.GetUserByID("u1")
	if err != nil || u.ID != "u1" {
		t.Errorf("expected user u1, got %v, err %v", u, err)
	}
}

func TestGetUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetUserByID("u1").Return(nil, errors.New("not found"))
	_, err := service.GetUserByID("u1")
	if err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
	if err := service.UpdateUser(&models.User{ID: "u1"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update error"))
	if err := service.UpdateUser(&models.User{ID: "u1"}); err == nil || err.Error() != "update error" {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	user := &models.User{Email: "john@example.com"}
	mockUserRepo.EXPECT().GetUserByEmail("john@example.com").Return(user, nil)
	u, err := service.GetUserByEmail(" john@example.com ")
	if err != nil || u.Email != "john@example.com" {
		t.Errorf("expected user john@example.com, got %v, err %v", u, err)
	}
}

func TestGetUserByEmail_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetUserByEmail("john@example.com").Return(nil, errors.New("not found"))
	_, err := service.GetUserByEmail(" john@example.com ")
	if err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	user := &models.User{Email: "john@example.com", Password: auth.HashPassword("pass")}
	mockUserRepo.EXPECT().GetUserByEmail("john@example.com").Return(user, nil)
	u, err := service.Login(" john@example.com ", "pass")
	if err != nil || u.Email != "john@example.com" {
		t.Errorf("expected login success, got %v, err %v", u, err)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	user := &models.User{Email: "john@example.com", Password: auth.HashPassword("pass")}
	mockUserRepo.EXPECT().GetUserByEmail("john@example.com").Return(user, nil)
	_, err := service.Login(" john@example.com ", "wrong")
	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected invalid credentials error, got %v", err)
	}
}

func TestLogin_GetUserByEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetUserByEmail("john@example.com").Return(nil, errors.New("not found"))
	_, err := service.Login(" john@example.com ", "pass")
	if err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestCreateUser_Fields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := setupUserService(ctrl)
	user := service.CreateUser(" John ", " john@example.com ", "pass", models.RoleGuest)
	if user.Name != "John" || user.Email != "john@example.com" || user.Role != models.RoleGuest || user.Available != false {
		t.Errorf("unexpected user fields: %+v", user)
	}
}

func TestCreateEmployee_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)
	mockUserRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
	emp, err := service.CreateEmployee("John", "john@example.com", "pass", models.RoleKitchenStaff, true)
	if err != nil || emp.Email != "john@example.com" {
		t.Errorf("expected employee creation, got %v, err %v", emp, err)
	}
}

func TestCreateEmployee_EmailInUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	users := []models.User{{Email: "john@example.com"}}
	mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
	_, err := service.CreateEmployee("John", "john@example.com", "pass", models.RoleKitchenStaff, true)
	if err == nil || err.Error() != "email already in use" {
		t.Errorf("expected email in use error, got %v", err)
	}
}

func TestCreateEmployee_SaveUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)
	mockUserRepo.EXPECT().SaveUser(gomock.Any()).Return(errors.New("save error"))
	_, err := service.CreateEmployee("John", "john@example.com", "pass", models.RoleKitchenStaff, true)
	if err == nil || err.Error() != "failed to save employee: save error" {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestGetUserNameByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	user := &models.User{ID: "u1", Name: "John"}
	mockUserRepo.EXPECT().GetUserByID("u1").Return(user, nil)
	name, err := service.GetUserNameByID("u1")
	if err != nil || name != "John" {
		t.Errorf("expected name John, got %v, err %v", name, err)
	}
}

func TestGetUserNameByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, mockUserRepo := setupUserService(ctrl)
	mockUserRepo.EXPECT().GetUserByID("u1").Return(nil, errors.New("not found"))
	_, err := service.GetUserNameByID("u1")
	if err == nil || err.Error() != "failed to fetch user: not found" {
		t.Errorf("expected fetch user error, got %v", err)
	}
}

package userService

import (
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/mock_userService.go -package=mocks

type IUserService interface {
	Signup(req dto.SignupRequest,role models.Role) error
	Login(email, password string) (*models.User, error)
	GetUserNameByID(userID string) (string, error)
	CreateUser(name, email, password string, role models.Role) models.User
	CreateEmployee(name, email, password string, role models.Role, available bool) (models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetAllEmployees() ([]models.User, error)
}

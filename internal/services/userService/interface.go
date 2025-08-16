package userService

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type IUserService interface {
	Signup(name, email, password string, roleint int) (string, error)
	Login(email, password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserNameByID(userID string) (string, error)
	CreateUser(name, email, password string, role models.Role) models.User
	CreateEmployee(name, email, password string, role models.Role, available bool) (models.User, error) 
}

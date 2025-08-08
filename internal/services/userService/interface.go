package userService

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type UserManager interface {
	Signup(name, email, password string, roleint int) (string, error)
	Login(email, password string) (*models.User, error)
	GetTotalGuests() (int, error)
	CreateUser(name, email, password string, role models.Role) models.User
	ReadPasswordMasked() (string, error)
}

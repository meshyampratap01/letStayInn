package userService

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
)

type UserManager interface {
	Signup(name, email, password string, roleint int) (string, error)
	Login(email, password string) (*models.User, error)
	GetTotalGuests() (int, error)
}

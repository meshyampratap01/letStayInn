package userRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type UserRepository interface {
	FindUserByEmail(user []models.User,email string) *models.User
	SaveUser(newUser models.User) error
	GetAllUsers() ([]models.User, error)
	SaveAllUsers([]models.User)error
}
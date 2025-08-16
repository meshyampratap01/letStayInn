package userRepository

import "github.com/meshyampratap01/letStayInn/internal/models"

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	FindUserByEmail(user []models.User, email string) *models.User
	SaveUser(newUser models.User) error
	GetAllUsers() ([]models.User, error)
	SaveAllUsers([]models.User)error
	GetUserByID(userID string) (*models.User, error)
	ToggleStaffAvailability(string) error
	GetStaffAvailability(string) (bool,error)
}
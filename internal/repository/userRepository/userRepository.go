package userRepository

import (
	"fmt"
	"strings"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/storage"
)

type FileUserRepository struct{}

func NewFileUserRepository() UserRepository {
	return &FileUserRepository{}
}

func (db *FileUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	if err := storage.ReadJson(config.UsersFile, &users); err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	return users, nil

}

func (db *FileUserRepository)FindUserByEmail(user []models.User, email string) *models.User {
	for _, u := range user {
		if u.Email == email {
			return &u
		}
	}
	return nil
}

func (db *FileUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var users []models.User
	if err := storage.ReadJson(config.UsersFile, &users); err != nil {
		return nil, err
	}

	for _, u := range users {
		if strings.EqualFold(u.Email, email) {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

func (db *FileUserRepository)SaveUser(newUser models.User) error {
	var users []models.User

	if err := storage.ReadJson(config.UsersFile, &users); err != nil {
		return fmt.Errorf("failed to read users: %w", err)
	}

	if db.FindUserByEmail(users, newUser.Email) != nil {
		return fmt.Errorf("email already exists")
	}

	users = append(users, newUser)

	if err := storage.WriteJson(config.UsersFile, users); err != nil {
		return fmt.Errorf("failed to write users: %w", err)
	}

	return nil
}

func (db *FileUserRepository) SaveAllUsers(users []models.User) error {
	if err := storage.WriteJson(config.UsersFile, users); err != nil {
		return fmt.Errorf("failed to write users: %w", err)
	}
	return nil
}

func (r *FileUserRepository) GetUserByID(userID string) (*models.User, error) {
	users, err := r.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	for _, user := range users {
		if userID == user.ID {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with ID %s not found", userID)
}

func (repo *FileUserRepository) ToggleStaffAvailability(userID string) error {
	var users []models.User
	err := storage.ReadJson(config.UsersFile, &users)
	if err != nil {
		return err
	}

	for i := range users {
		if users[i].ID == userID {
			users[i].Available = !users[i].Available
			return storage.WriteJson(config.UsersFile, users)
		}
	}
	return fmt.Errorf("user not found")
}

func (repo *FileUserRepository) GetStaffAvailability(userID string) (bool, error) {
	var users []models.User
	err := storage.ReadJson(config.UsersFile, &users)
	if err != nil {
		return false, err
	}

	for _, u := range users {
		if u.ID == userID {
			return u.Available, nil
		}
	}
	return false, fmt.Errorf("user not found")
}

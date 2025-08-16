package userService

import (
	"fmt"
	"strings"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
)

type UserService struct {
	userRepo userRepository.UserRepository
}

func NewUserService(userRepo userRepository.UserRepository) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Signup(name, email, password string, roleint int) (string, error) {
	role := models.Role(roleint)
	if role < models.RoleGuest || role > models.RoleManager {
		return "", fmt.Errorf("invalid role")
	}

	newUser := s.CreateUser(name, email, password, role)

	if err := s.userRepo.SaveUser(newUser); err != nil {
		return "", err
	}

	return "Signup successful as Guest!! Please login.", nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	email = strings.TrimSpace(email)
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (us *UserService) CreateUser(name, email, password string, role models.Role) models.User {
	return models.User{
		ID:        utils.NewUUID(),
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		Password:  auth.HashPassword(password),
		Role:      role,
		CreatedAt: time.Now(),
		Available: role != models.RoleGuest,
	}
}

func (s *UserService) CreateEmployee(name, email, password string, role models.Role, available bool) (models.User, error) {
	if role != models.RoleKitchenStaff && role != models.RoleCleaningStaff && role != models.RoleManager {
		return models.User{}, fmt.Errorf("invalid role for employee")
	}

	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to fetch users: %v", err)
	}
	for _, u := range users {
		if u.Email == email {
			return models.User{}, fmt.Errorf("email already in use")
		}
	}

	newUser := models.User{
		ID:        utils.NewUUID(),
		Name:      name,
		Email:     email,
		Password:  auth.HashPassword(password),
		Role:      role,
		CreatedAt: time.Now(),
		Available: available,
	}

	if err := s.userRepo.SaveUser(newUser); err != nil {
		return models.User{}, fmt.Errorf("failed to save employee: %v", err)
	}

	return newUser, nil
}

func (s *UserService) GetUserNameByID(userID string) (string, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	return user.Name, nil
}

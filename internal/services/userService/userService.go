package userService

import (
	"fmt"
	mail "net/mail"
	"strings"

	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/utils"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type UserService struct {
	userRepo userRepository.UserRepository
}

func NewUserService(userRepo userRepository.UserRepository) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Signup(req dto.SignupRequest, role models.Role) error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf(constants.ErrInvalidCredentials)
	}

	if err := validators.ValidatePassword(req.Password); err != nil {
		return fmt.Errorf(constants.ErrInvalidCredentials)
	}

	newUser := s.CreateUser(req.Name, req.Email, req.Password, role)

	if err := s.userRepo.SaveUser(newUser); err != nil {
		return fmt.Errorf(constants.ErrFailedToSaveUser)
	}

	return nil
}

func (s *UserService) GetAllEmployees() ([]models.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var employees []models.User
	for _, u := range users {
		if u.Role == models.RoleKitchenStaff || u.Role == models.RoleCleaningStaff || u.Role == models.RoleManager {
			employees = append(employees, u)
		}
	}
	return employees, nil
}
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}


func (s *UserService) Login(email, password string) (*models.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidCredentials)
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf(constants.ErrInvalidCredentials)
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

package userService

import (
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/repository/userRepository"
	"github.com/meshyampratap01/letStayInn/internal/storage"
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

func (s *UserService) Login(email, password string) (*models.User, error) {
	var users []models.User
	storage.ReadJson(config.UsersFile, &users)

	user := s.userRepo.FindUserByEmail(users, strings.TrimSpace(email))
	if user == nil || !auth.CheckPassword(user.Password, strings.TrimSpace(password)) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (us *UserService) GetTotalGuests() (int, error) {
	users, err := us.userRepo.GetAllUsers()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, u := range users {
		if u.Role == models.RoleGuest {
			count++
		}
	}
	return count, nil
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

func (us *UserService) ReadPasswordMasked() (string, error) {
	if err := keyboard.Open(); err != nil {
		return "", err
	}
	defer keyboard.Close()
	password := ""
	for {
		char, key, err := keyboard.GetKey()
		if err!=nil{
			return "",err
		}
		if key == keyboard.KeyEnter{
			break
		} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2{
			if len(password)>0{
				password = password[:len(password)-1]
				fmt.Print("\b \b")
			}
		} else{
			password +=string(char)
			fmt.Print("*")
		}
	}
	fmt.Println()
	return password,nil
}

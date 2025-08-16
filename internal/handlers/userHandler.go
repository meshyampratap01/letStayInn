package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"github.com/meshyampratap01/letStayInn/internal/utils"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type UserHandler struct {
	userService      userService.IUserService
	DashboardHandler *DashboardHandler
	feedbackService  feedbackService.IFeedbackService
}

func NewUserHandler(us userService.IUserService, DashboardHandler *DashboardHandler, fs feedbackService.IFeedbackService) *UserHandler {
	return &UserHandler{
		userService:      us,
		DashboardHandler: DashboardHandler,
		feedbackService:  fs,
	}
}

func (u *UserHandler) SignupHandler() {
	reader := bufio.NewReader(os.Stdin)

	color.Cyan(config.SignupMsg)
	fmt.Print(color.HiWhiteString("Enter name: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var email string
	for {
		fmt.Print(color.HiWhiteString("Enter email: "))
		emailInput, _ := reader.ReadString('\n')
		email = strings.TrimSpace(emailInput)

		if err := validators.ValidateEmail(email); err != nil {
			color.Red("Error: %v", err)
			continue
		}
		break
	}

	var password string
	for {
		pass, err := utils.ReadPasswordMasked(color.HiWhiteString("Enter password: "))
		if err != nil {
			color.Red("Error reading password: %v", err)
			continue
		}
		password = pass
		if err := validators.ValidatePassword(password); err != nil {
			color.Red("Error: %v", err)
			continue
		}
		break
	}

	roleint := 1

	msg, err := u.userService.Signup(name, email, password, roleint)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	color.Green(msg)
}

func (u *UserHandler) LoginHandler() {
	reader := bufio.NewReader(os.Stdin)

	color.Cyan(config.LoginMsg)
	fmt.Print(color.HiWhiteString("Enter email: "))
	email, _ := reader.ReadString('\n')

	password, err := utils.ReadPasswordMasked(color.HiWhiteString("Enter password: "))
	if err != nil {
		color.Red("Error reading password: %v", err)
		return
	}

	user, err := u.userService.Login(email, password)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, contextkeys.UserRoleKey, user.Role)
	ctx = context.WithValue(ctx, contextkeys.UserNameKey, user.Name)

	color.Green(config.UserWelcome+"%s!", user.Name)
	u.DashboardHandler.LoadDashboard(ctx)
}

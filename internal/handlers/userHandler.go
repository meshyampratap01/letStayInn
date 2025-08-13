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
		fmt.Print(color.HiWhiteString("Enter Password: "))

		pass,err:=u.userService.ReadPasswordMasked()
		if err!=nil{
			color.Red("Error reading password: %v",err)
			continue
		}
		password = strings.TrimSpace(pass)
		if err:=validators.ValidatePassword(password); err!=nil{
			color.Red("Error: %v",err)
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

	fmt.Print(color.HiWhiteString("Enter password: "))
	password, _ := reader.ReadString('\n')

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

func (h *UserHandler) SubmitFeedback(ctx context.Context) {
	if name, ok := ctx.Value(contextkeys.UserNameKey).(string); ok {
		color.Cyan(config.FeedbackMsg, name)
	}

	err := h.feedbackService.SubmitFeedback(ctx)
	if err != nil {
		color.Red("Error submitting feedback: %v", err)
	}
}

package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/services/feedbackService"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type UserHandler struct {
	userService      userService.UserManager
	DashboardHandler *DashboardHandler
	feedbackService  feedbackService.FeedbackServiceManager
}

func NewUserHandler(us userService.UserManager, DashboardHandler *DashboardHandler,fs feedbackService.FeedbackServiceManager) *UserHandler {
	return &UserHandler{
		userService:      us,
		DashboardHandler: DashboardHandler,
		feedbackService: fs,
	}
}

func (u *UserHandler) SignupHandler() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n---- SignUp ---- ")
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var email string
	for {
		fmt.Print("Enter email: ")
		emailInput, _ := reader.ReadString('\n')
		email = strings.TrimSpace(emailInput)

		if err := validators.ValidateEmail(email); err != nil {
			fmt.Println("Error:", err)
			continue
		}
		break
	}

	var password string
	for {
		fmt.Print("Enter password: ")
		passwordInput, _ := reader.ReadString('\n')
		password = strings.TrimSpace(passwordInput)

		if err := validators.ValidatePassword(password); err != nil {
			fmt.Println("Error:", err)
			continue
		}
		break
	}

	roleint := 1 

	msg, err := u.userService.Signup(name, email, password, roleint)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(msg)
}


func (u *UserHandler) LoginHandler() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n---- Login ---- ")
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')

	user, err := u.userService.Login(email, password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, contextkeys.UserRoleKey, user.Role)
	ctx = context.WithValue(ctx, contextkeys.UserNameKey, user.Name)

	fmt.Printf("Welcome, %s!\n", user.Name)
	u.DashboardHandler.LoadDashboard(ctx)
}

func (h *UserHandler) SubmitFeedback(ctx context.Context) {
	if name, ok := ctx.Value(contextkeys.UserNameKey).(string); ok {
		fmt.Printf("Hello %s! We'd love to hear your thoughts.\n", name)
	}

	err := h.feedbackService.SubmitFeedback(ctx)
	if err != nil {
		fmt.Println("Error submitting feedback:", err)
	}
}
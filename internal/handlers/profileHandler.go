package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/auth"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/services/userService"
	"github.com/meshyampratap01/letStayInn/internal/validators"
)

type ProfileHandler struct {
	userService userService.IUserService
}

func NewProfileHandler(userService userService.IUserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (h *ProfileHandler) ViewProfile(ctx context.Context) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		color.Red("User ID not found in context.")
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		color.Red("Error fetching user profile: %v", err)
		return
	}
	color.Cyan("\n--- Your Profile ---")
	fmt.Printf("Name: %s\nEmail: %s\nRole: %s\n", user.Name, user.Email, user.Role.String())
}

func (h *ProfileHandler) UpdateProfile(ctx context.Context) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok {
		color.Red("User ID not found in context.")
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		color.Red("Error fetching user profile: %v", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("\n--- Update Profile ---")
	updated := false
	for {
		fmt.Println("What would you like to update?")
		fmt.Println("1. Name")
		fmt.Println("2. Email")
		fmt.Println("3. Password")
		fmt.Println("4. Done/Cancel")
		fmt.Print("Enter choice: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		switch choiceStr {
		case "1":
			fmt.Printf("Current Name: %s\n", user.Name)
			fmt.Print("Enter new name (leave blank to keep current): ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name != "" {
				user.Name = name
				updated = true
				color.Green("Name will be updated.")
			} else {
				color.Yellow("Name unchanged.")
			}
		case "2":
			fmt.Printf("Current Email: %s\n", user.Email)
			fmt.Print("Enter new email (leave blank to keep current): ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)
			if email != "" {
				user.Email = email
				updated = true
				color.Green("Email will be updated.")
			} else {
				color.Yellow("Email unchanged.")
			}
		case "3":
			fmt.Print("Enter current password: ")
			currPwd, _ := reader.ReadString('\n')
			currPwd = strings.TrimSpace(currPwd)
			if !auth.CheckPassword(user.Password, currPwd) {
				color.Red("Current password is incorrect. Password not updated.")
			} else {
				fmt.Print("Enter new password: ")
				newPwd, _ := reader.ReadString('\n')
				newPwd = strings.TrimSpace(newPwd)
				if newPwd == "" {
					color.Red("New password cannot be empty. Password not updated.")
				} else if err := validators.ValidatePassword(newPwd); err != nil {
					color.Red("%v. Password not updated.", err)
				} else {
					user.Password = auth.HashPassword(newPwd)
					updated = true
					color.Green("Password will be updated.")
				}
			}
		case "4":
			if updated {
				if err := h.userService.UpdateUser(user); err != nil {
					color.Red("Error updating profile: %v", err)
					return
				}
				color.Green("Profile updated successfully!")
			} else {
				color.Yellow("No changes made to profile.")
			}
			return
		default:
			color.Red("Invalid choice. Please select a valid option.")
		}
	}
}

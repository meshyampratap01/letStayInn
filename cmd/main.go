package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/container"
)

func main() {
	CLIUserHandler := container.InitHandlers()

	for {
		color.Cyan(config.WelcomeMsg)
		color.Cyan(" "+config.AppDescription+" ")
		color.Yellow("1. Signup")
		color.Yellow("2. Login")
		color.Yellow("3. Exit")

		fmt.Print(color.HiWhiteString("Select Option: "))
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CLIUserHandler.SignupHandler()
		case 2:
			CLIUserHandler.LoginHandler()
		case 3:
			color.Green("Exiting...")
			os.Exit(0)
		default:
			color.Red(config.InvalidOption)
		}
	}
}

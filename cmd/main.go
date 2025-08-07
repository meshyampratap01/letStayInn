package main

import (
	"fmt"
	"os"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/container" 
)

func main() {
	CLIUserHandler := container.InitHandlers()

	for {
		fmt.Println(config.LoginMsg)
		fmt.Println("1.Signup")
		fmt.Println("2.Login")
		fmt.Println("3.Exit")
		fmt.Print("Select Option: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CLIUserHandler.SignupHandler()
		case 2:
			CLIUserHandler.LoginHandler()
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println(config.InvalidOption)
		}
	}
}

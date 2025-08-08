package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AddBackButton() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter 'b' to go back: ")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "b" {
			break
		}
		fmt.Print("Invalid input. Please enter 'b' to go back: ")
	}
}

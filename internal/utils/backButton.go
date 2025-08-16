package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AddBackButton() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress Enter to go back...")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" { 
			break
		}
		fmt.Print("Invalid input. Just press Enter to go back: ")
	}
}


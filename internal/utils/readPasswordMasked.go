package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ReadPasswordMasked(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePassword)), nil
}

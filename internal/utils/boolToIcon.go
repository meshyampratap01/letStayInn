package utils

import "github.com/fatih/color"

func BoolToIcon(value bool, yesText, noText string) string {
	if value {
		return color.GreenString("✔ " + yesText)
	}
	return color.RedString("✘ " + noText)
}

package validators

import (
	"errors"
	"regexp"
	"time"
)

func ValidateEmail(email string) error {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !regex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}



func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	letterRegex := regexp.MustCompile(`[A-Za-z]`)
	numberRegex := regexp.MustCompile(`\d`)
	specialCharRegex := regexp.MustCompile(`[@$!%*#?&]`)

	if !letterRegex.MatchString(password) {
		return errors.New("password must include at least one letter")
	}
	if !numberRegex.MatchString(password) {
		return errors.New("password must include at least one number")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must include at least one special character (@$!%*#?&)")
	}
	return nil
}


func ValidateDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format (expected YYYY-MM-DD)")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return time.Time{}, errors.New("date must be today or in the future")
	}

	return parsedDate, nil
}

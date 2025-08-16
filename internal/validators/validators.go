package validators

import (
	"errors"
	"regexp"
	"time"
)

func ValidateEmail(email string) error {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !regex.MatchString(email) {
		return errors.New("invalid email format, enter a valid email(eg. sh.singh@gmail.com)")
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
		return errors.New("password must alphanumeric with atleast one special character (@$!%*#?&)")
	}
	if !numberRegex.MatchString(password) {
		return errors.New("password must alphanumeric with atleast one special character (@$!%*#?&)")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must alphanumeric with atleast one special character (@$!%*#?&)")
	}
	return nil
}


func ValidateDate(dateStr string) (string, error) {
	layout := "02-01-2006"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", errors.New("invalid date format (expected DD-MM-YYYY)")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		return "", errors.New("date must be today or in the future")
	}

	return parsedDate.Format(layout), nil
}


func ValidateCheckoutDate(checkinStr, checkoutStr string) (string, error) {
	layout := "02-01-2006"

	checkinDate, err := time.Parse(layout, checkinStr)
	if err != nil {
		return "", errors.New("invalid check-in date format (expected DD-MM-YYYY)")
	}

	checkoutDate, err := time.Parse(layout, checkoutStr)
	if err != nil {
		return "", errors.New("invalid check-out date format (expected DD-MM-YYYY)")
	}

	if checkoutDate.Before(checkinDate) {
		return "", errors.New("checkout date must be the same or after the check-in date")
	}

	return checkoutDate.Format(layout), nil
}




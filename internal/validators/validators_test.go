package validators

import (
	"testing"
	"time"
)


func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"user@example.com", false},
		{"invalid-email", true},
		{"user@com", true},
		{"", true},
		{"user@domain.co.in", false},
	}

	for _, test := range tests {
		err := ValidateEmail(test.email)
		if (err != nil) != test.wantErr {
			t.Errorf("ValidateEmail(%q) = %v, wantErr %v", test.email, err, test.wantErr)
		}
	}
}


func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"Pass123!", false},
		{"weak", true},
		{"12345678", true},
		{"Password", true},
		{"Pass123", true},
		{"Valid$123", false},
	}

	for _, test := range tests {
		err := ValidatePassword(test.password)
		if (err != nil) != test.wantErr {
			t.Errorf("ValidatePassword(%q) = %v, wantErr %v", test.password, err, test.wantErr)
		}
	}
}


func TestValidateDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	future := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	tests := []struct {
		dateStr string
		wantErr bool
	}{
		{today, false},
		{yesterday, true},
		{future, false},
		{"2021-02-30", true}, // invalid date
		{"not-a-date", true},
	}

	for _, test := range tests {
		_, err := ValidateDate(test.dateStr)
		if (err != nil) != test.wantErr {
			t.Errorf("ValidateDate(%q) = %v, wantErr %v", test.dateStr, err, test.wantErr)
		}
	}
}

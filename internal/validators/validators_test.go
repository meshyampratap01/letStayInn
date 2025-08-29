package validators

import (
	"testing"
	"time"

	"github.com/meshyampratap01/letStayInn/internal/models"
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
	today := time.Now().Format("02-01-2006")
	future := time.Now().Add(24 * time.Hour).Format("02-01-2006")
	past := time.Now().Add(-24 * time.Hour).Format("02-01-2006")

	tests := []struct {
		dateStr string
		wantErr bool
	}{
		{today, false},       
		{future, false},      
		{past, true},         
		{"31-02-2024", true}, 
		{"2024-12-31", true}, 
	}

	for _, test := range tests {
		_, err := ValidateDate(test.dateStr)
		if (err != nil) != test.wantErr {
			t.Errorf("ValidateDate(%q) error = %v, wantErr %v", test.dateStr, err, test.wantErr)
		}
	}
}

func TestValidateCheckoutDate(t *testing.T) {
	layout := "02-01-2006"
	checkin := time.Now().Format(layout)
	checkoutAfter := time.Now().Add(24 * time.Hour).Format(layout)
	checkoutSame := checkin
	checkoutBefore := time.Now().Add(-24 * time.Hour).Format(layout)

	tests := []struct {
		checkin  string
		checkout string
		wantErr  bool
	}{
		{checkin, checkoutAfter, false},  
		{checkin, checkoutSame, false},   
		{checkin, checkoutBefore, true},  
		{"invalid", checkoutAfter, true}, 
		{checkin, "invalid", true},       
	}

	for _, test := range tests {
		_, err := ValidateCheckoutDate(test.checkin, test.checkout)
		if (err != nil) != test.wantErr {
			t.Errorf("ValidateCheckoutDate(%q, %q) error = %v, wantErr %v", test.checkin, test.checkout, err, test.wantErr)
		}
	}
}

func TestIsValidRoomType(t *testing.T) {
	validTypes := []models.RoomType{
		models.RoomTypeStandard,
		models.RoomTypeDeluxe,
		models.RoomTypeSuite,
		models.RoomTypeExecutive,
	}
	for _, rt := range validTypes {
		if !IsValidRoomType(string(rt)) {
			t.Errorf("IsValidRoomType(%q) = false, want true", rt)
		}
	}

	if IsValidRoomType("Penthouse") {
		t.Errorf("IsValidRoomType(%q) = true, want false", "Penthouse")
	}
}

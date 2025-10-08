package auth

import (
	"testing"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "mySecret123"
	hashed := HashPassword(password)
	if hashed == "" {
		t.Error("Expected non-empty hashed password")
	}
	if !CheckPassword(hashed, password) {
		t.Error("CheckPassword should return true for correct password")
	}
	if CheckPassword(hashed, "wrongPassword") {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestHashAndCheck(t *testing.T) {
	hash, err := HashPassword("Correct-Horse-9!")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "Correct-Horse-9!") {
		t.Error("correct password was rejected")
	}
	if CheckPassword(hash, "Wrong-Horse-9!") {
		t.Error("wrong password was accepted")
	}
}

// internal/auth/password_test.go (add to the file)
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		pw      string
		wantErr error
	}{
		{"Abcdef1!", nil},
		{"Ab1!", ErrPasswordTooShort},
		{"abcdefg1!", ErrWeakPassword},                   // no uppercase
		{"ABCDEFG1!", ErrWeakPassword},                   // no lowercase
		{"Abcdefgh!", ErrWeakPassword},                   // no digit
		{"Abcdefg12", ErrWeakPassword},                   // no symbol
		{strings.Repeat("Aa1!", 19), ErrPasswordTooLong}, // 76 bytes
	}

	for _, tc := range tests {
		if err := ValidatePassword(tc.pw); !errors.Is(err, tc.wantErr) {
			t.Errorf("ValidatePassword(%q) = %v, want %v", tc.pw, err, tc.wantErr)
		}
	}
}

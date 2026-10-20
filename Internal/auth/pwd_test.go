package auth

import "testing"

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

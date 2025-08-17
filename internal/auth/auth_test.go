package auth

import (
	"testing"
)

func TestHashPassword_EmptyPassword(t *testing.T) {
	hashPass, err := HashPassword("")

	if hashPass != "" || err == nil {
		t.Errorf(`HashPassowrd("") = %q, %v, want "", error`, hashPass, err)
	}

}

func TestHashPassword_WithPassword(t *testing.T) {
	password := "validpassword123"
	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword() error = %v, want nil", err)
	}

	if got, want := hash, ""; got == want {
		t.Errorf("HashPassowrd() = %q, want non-empty string", got)
	}
}

func TestHashPassword_WithReversePassword(t *testing.T) {
	password := "hello123!"
	reversepassword := "!321olleh"
	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword() error = %v, want nil", err)
	}

	reversePasshash, err := HashPassword(reversepassword)
	if err != nil {
		t.Errorf("HashPassword() error = %v, want nil", err)
	}

	if reversePasshash == hash {
		t.Errorf("HashPassword() hash = %v == reversePasshash = %v, want different hash", hash, reversePasshash)
	}
}

func TestHashPassword_SamePasswordDifferentHashes(t *testing.T) {
	password := "samepassword"
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("Same password should produce different hashes (salt)")
	}
}

func TestHashPassword_SimilarPasswords(t *testing.T) {
	hash1, _ := HashPassword("password1")
	hash2, _ := HashPassword("password2")

	if hash1 == hash2 {
		t.Error("Similar passwords should produce different hashes")
	}
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	password := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	_, err := HashPassword(password)

	if err != nil {
		t.Error("Special characters password should produce hash")
	}
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	hash := "$2a$10$validhashexample"
	err := CheckPasswordHash("", hash)

	if err == nil {
		t.Error("CheckPasswordHash(\"\", hash) should return error for empty password")
	}
}

func TestCheckPasswordHash_ValidPassword(t *testing.T) {
	password := "validpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("CheckPasswordHash() error = %v, want nil", err)
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "correctpassword"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Error("CheckPasswordHash() should return error for wrong password")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	password := "validpassword"
	invalidHash := "notavalidhash"

	err := CheckPasswordHash(password, invalidHash)
	if err == nil {
		t.Error("CheckPasswordHash() should return error for invalid hash")
	}
}

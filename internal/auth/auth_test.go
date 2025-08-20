package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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

// Test cases for MakeJWT function
func TestMakeJWT_ValidInput(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiration := time.Hour

	token, err := MakeJWT(userID, secret, expiration)

	if err != nil {
		t.Errorf("MakeJWT() error = %v, want nil", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}

	// Verify we can validate the token we just created
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Failed to validate generated token: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("ValidateJWT() userID = %v, want %v", parsedUserID, userID)
	}
}

func TestMakeJWT_EmptySecret(t *testing.T) {
	userID := uuid.New()
	secret := ""
	expiration := time.Hour

	token, err := MakeJWT(userID, secret, expiration)

	// Token should still be created even with empty secret
	if err != nil {
		t.Errorf("MakeJWT() error = %v, want nil", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}
}

// Test cases for ValidateJWT function
func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiration := time.Hour

	// First create a valid token
	token, err := MakeJWT(userID, secret, expiration)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Now validate it
	parsedUserID, err := ValidateJWT(token, secret)

	if err != nil {
		t.Errorf("ValidateJWT() error = %v, want nil", err)
	}

	if parsedUserID != userID {
		t.Errorf("ValidateJWT() userID = %v, want %v", parsedUserID, userID)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "test-secret-key"
	invalidToken := "invalid.jwt.token"

	userID, err := ValidateJWT(invalidToken, secret)

	if err == nil {
		t.Error("ValidateJWT() should return error for invalid token")
	}

	if userID != uuid.Nil {
		t.Errorf("ValidateJWT() should return uuid.Nil for invalid token, got %v", userID)
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer token",
			authHeader:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError:   false,
		},
		{
			name:          "Valid Bearer token with extra spaces",
			authHeader:    "Bearer   token123   ",
			expectedToken: "token123",
			expectError:   false,
		},
		{
			name:          "Missing Authorization header",
			authHeader:    "",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "No Bearer prefix",
			authHeader:    "Basic dXNlcjpwYXNz",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Bearer without space",
			authHeader:    "Bearertoken123",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Bearer with no token",
			authHeader:    "Bearer",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Bearer with only spaces",
			authHeader:    "Bearer   ",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Wrong prefix entirely",
			authHeader:    "NotBearer token123",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Case sensitive Bearer",
			authHeader:    "bearer token123",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Multiple spaces between Bearer and token",
			authHeader:    "Bearer     mytoken",
			expectedToken: "mytoken",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create headers
			headers := http.Header{}
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			// Call the function
			token, err := GetBearerToken(headers)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check token value
			if token != tt.expectedToken {
				t.Errorf("Expected token %q but got %q", tt.expectedToken, token)
			}
		})
	}
}

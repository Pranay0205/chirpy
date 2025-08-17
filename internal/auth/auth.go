package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil

}

func CheckPasswordHash(password, hash string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}

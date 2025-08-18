package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

type CustomClaim struct {
	jwt.RegisteredClaims
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claim := CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "chirpy",
			Subject:   userId.String(),
		},
	}

	signingMethod := jwt.SigningMethodHS256

	token := jwt.NewWithClaims(signingMethod, claim)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("failed to validate token: algorithm change detected")
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed while validating token: %w", err)
	}
	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to retreive userId while validating token: %w", err)
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse userId while validating token: %w", err)
	}

	return userId, nil
}

package main

import (
	"time"

	"github.com/google/uuid"
)

// Response Models
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

// Request Models
type ChirpRequest struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type UserRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ExpiresInSec int64  `json:"expires_in_seconds"`
}

package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userReq := UserRequest{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse the request: invalid JSON", err)
		return
	}
	hashed_password, err := auth.HashPassword(userReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create password hash: incompatible password format", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{Email: userReq.Email, HashedPassword: hashed_password})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user: database error", err)
		return
	}

	user := User{
		ID:         dbUser.ID,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
		Email:      dbUser.Email,
		IsChirpRed: dbUser.IsChirpRed.Bool,
	}

	respondWithJSON(w, 201, user)
}

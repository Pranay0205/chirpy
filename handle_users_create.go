package main

import (
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

	dbUser, err := cfg.dbQueries.CreateUser(context.Background(), userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user: database error", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, user)
}

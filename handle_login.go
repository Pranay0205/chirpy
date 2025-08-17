package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userCredentials := UserRequest{}
	err := decoder.Decode(&userCredentials)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse login request: invalid JSON", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), userCredentials.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "failed to login: incorrect email or password", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "failed to get user: database error", err)
		return
	}

	err = auth.CheckPasswordHash(userCredentials.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to login: incorrect email or password", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}

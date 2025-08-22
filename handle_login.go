package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
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

		respondWithError(w, http.StatusInternalServerError, "failed to login: database error", err)
		return
	}

	err = auth.CheckPasswordHash(userCredentials.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to login: incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to login: token creationg failed", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to login: unable to generate refresh token", err)
		return
	}

	dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
		UpdatedAt: time.Now().UTC(),
		UserID:    dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to login: database error", err)
		return
	}

	user := User{
		ID:            dbUser.ID,
		CreatedAt:     dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
		Email:         dbUser.Email,
		Token:         token,
		Refresh_Token: dbRefreshToken.Token,
		IsChirpRed:    dbUser.IsChirpRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)
}

package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleUpdateUsers(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to update user: unauthorized user", err)
		return
	}

	tokenUserId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to update user: unauthorized user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	userReq := UserRequest{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse request: invalid JSON", err)
		return
	}

	dbUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:    tokenUserId,
		Email: userReq.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user: database error", err)
		return
	}

	user := User{
		ID:         dbUser.ID,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
		Email:      dbUser.Email,
		IsChirpRed: dbUser.IsChirpRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)
}

package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirps(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to delete the chirp: unauthorized user", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to delete the chirp: unauthorized user", err)
		return
	}

	chirpIDValue := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse chirp ID: invalid UUID format", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "failed to delete chirp: chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to delete the chirp: database error", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "failed to delete chirp: you can only delete your own chirps", nil)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete the chirp: database error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

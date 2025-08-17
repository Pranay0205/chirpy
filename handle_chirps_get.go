package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, r *http.Request) {

	chirpIDValue := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse chirp ID: invalid UUID format", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "failed to get chirp: chirp not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "failed to get chirp: database error", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}

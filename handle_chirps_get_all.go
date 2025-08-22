package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {

	authorIdString := r.URL.Query().Get("author_id")

	authorId, err := uuid.Parse(authorIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to retreive all chirps: invalid author_id", err)
		return
	}

	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context(), authorId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retreive all chirps: database error", err)
		return
	}

	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{ID: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: dbChirp.UserID})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

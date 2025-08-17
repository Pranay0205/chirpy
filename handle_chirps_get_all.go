package main

import (
	"net/http"
)

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
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

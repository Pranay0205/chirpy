package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	requestVal := ChirpRequest{}
	err := decoder.Decode(&requestVal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse request: invalid JSON", err)
		return
	}

	const maxChirpLength = 140
	if len(requestVal.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "failed to create chirp: body exceeds 140 characters", nil)
		return
	}

	censoredChirp := censorProfaneWords(requestVal.Body)
	parsedUUID, err := uuid.Parse(requestVal.UserId)
	if err != nil {
		log.Printf("Error parsing UUID string: %v\n", err)
		respondWithError(w, http.StatusBadRequest, "failed to parse user ID: invalid UUID format", err)
		return
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredChirp,
		UserID: parsedUUID,
	})
	if err != nil {
		log.Printf("Error while creating the chirp: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp: database error", err)
		return
	}

	chirp := Chirp{ID: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: dbChirp.UserID}

	respondWithJSON(w, 201, chirp)

}

func censorProfaneWords(chirp string) string {

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	chirpWords := strings.Split(chirp, " ")
	for i, word := range chirpWords {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			chirpWords[i] = "****"
		}
	}

	return strings.Join(chirpWords, " ")
}

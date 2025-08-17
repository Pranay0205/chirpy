package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	requestVal := parameters{}
	err := decoder.Decode(&requestVal)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(requestVal.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	censoredChirp := censorProfaneWords(requestVal.Body)
	parsedUUID, err := uuid.Parse(requestVal.UserId)
	if err != nil {
		log.Printf("Error parsing UUID string: %v\n", err)
		respondWithError(w, http.StatusBadRequest, "Error parsing UUID string", err)
		return
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredChirp,
		UserID: parsedUUID,
	})
	if err != nil {
		log.Printf("Error while creating the chirp: %v\n", err)
		respondWithError(w, http.StatusBadRequest, "Error while creating the chirp", err)
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

package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to create chirp: unauthorized user", err)
		return
	}

	tokenUserId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to create chirp: unauthorized user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	requestVal := ChirpRequest{}
	err = decoder.Decode(&requestVal)
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

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredChirp,
		UserID: tokenUserId,
	})
	if err != nil {
		log.Printf("Error while creating the chirp: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp: database error", err)
		return
	}

	chirp := Chirp{ID: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: tokenUserId}

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

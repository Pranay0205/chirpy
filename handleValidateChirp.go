package main

import (
	"encoding/json"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		Valid bool `json:"valid"`
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

	respondWithJSON(w, http.StatusOK, returnValue{Valid: true})

}

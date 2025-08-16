package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		CleanedBody string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnValue{CleanedBody: censoredChirp})

}

func censorProfaneWords(chirp string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	chirpWords := strings.Split(chirp, " ")
	for i, word := range chirpWords {

		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				chirpWords[i] = "****"
			}
		}
	}

	cleanString := strings.Join(chirpWords, " ")

	return cleanString
}

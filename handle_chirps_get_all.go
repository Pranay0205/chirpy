package main

import (
	"chirpy/internal/database"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve all chirps: database error", err)
		return
	}

	if authorIdString != "" {
		authorId, err := uuid.Parse(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "failed to retrieve all chirps: invalid author_id", err)
			return
		}

		var filteredChirps []database.Chirp
		for _, chirp := range dbChirps {
			if chirp.UserID == authorId {
				filteredChirps = append(filteredChirps, chirp)
			}
		}
		dbChirps = filteredChirps
	}

	if sortParam == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
	}

	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{ID: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: dbChirp.UserID})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

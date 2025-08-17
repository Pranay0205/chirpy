package main

import (
	"context"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "failed to reset metrics: not dev environment", errors.New("forbidden in production"))
		return
	}

	cfg.fileserverHits.Store(0)
	cfg.dbQueries.DeleteAllUsers(context.Background())
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hits reset to 0"))
}

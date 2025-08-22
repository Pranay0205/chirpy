package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	webHookReq := WebHooksRequest{}
	err = decoder.Decode(&webHookReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode webhook payload: bad request", err)
		return
	}

	if webHookReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(webHookReq.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to process request: invalid user ID", err)
		return
	}

	err = cfg.dbQueries.UpgradeUserToChirpRed(r.Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "failed to upgrade user: user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to upgrade user: database error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

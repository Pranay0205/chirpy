package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to refresh token: no token found", err)
		return
	}

	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to refresh token: database error", err)
		return
	}

	if time.Now().UTC().After(dbRefreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "failed to refresh token: refresh token expired", nil)
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "failed to refresh token: refresh token revoked", nil)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to refresh token: database error", err)
		return
	}

	newToken, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to refresh token: new token generation failed", err)
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{Token: newToken})

}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid request", err)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		Token:     token,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to revoke token: database error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

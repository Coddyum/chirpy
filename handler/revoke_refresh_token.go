package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
)

func (cfg *ApiConfig) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Impossible de récupérer le token %s", err)
		return
	}

	existingToken, err := cfg.DB.GetRefreshTokenInfo(r.Context(), token)
	if err != nil {
		log.Printf("Impossible de récupérer le token ou alors le token n'existe pas dans la db: %s", err)
		return
	}

	err = cfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:     existingToken.Token,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		log.Printf("Impossible d'update la base de donner refreshToken: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

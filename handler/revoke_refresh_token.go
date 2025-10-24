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
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Impossible de récupérer le refreshToken dans le header: %s", err)
		return
	}

	existingToken, err := cfg.DB.GetRefreshTokenInfo(r.Context(), headerToken)
	if err != nil {
		log.Printf("Impossible de récupérer les informations du token: %s", err)
		return
	}

	if existingToken.Token != headerToken {
		log.Println("Le token passer dans le header n'est pas equale a celui de la db")
		return
	}

	now := time.Now()

	err = cfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:     headerToken,
		UpdatedAt: now,
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	})
	if err != nil {
		log.Println("Impossible d'update la db")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

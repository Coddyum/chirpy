package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/utils"
)

type NewToken struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Impossible de récupérer le refreshToken dans le header %s", err)
		return
	}

	checkTokenInfo, err := cfg.DB.GetRefreshTokenInfo(r.Context(), token)
	if err != nil {
		log.Printf("Impossible de récupérer les informations sur le token %s", err)
		return
	}

	if checkTokenInfo.ExpiresAt.Time.Before(time.Now().UTC()) || checkTokenInfo.RevokedAt.Valid {
		log.Printf("Le token a expirer le: %s ou a été revoked le %s ", checkTokenInfo.ExpiresAt.Time, checkTokenInfo.RevokedAt.Time)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := checkTokenInfo.UserID
	newAccessToken, err := auth.MakeJWT(userID.UUID, cfg.JWTSecret)
	if err != nil {
		log.Printf("Erreur lors de la création du nouveau JWT access token: %s", err)
		return
	}

	utils.WriteJson(w, 200, Refresh{Token: newAccessToken})
}

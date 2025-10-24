package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/utils"
)

type NewToken struct {
	AccessToken string `json:"access_token"`
}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Impossible de récupèrer le bearer Token: %s", err)
		return
	}

	// Check if the token Exist and return The Users
	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), headerToken)
	if err != nil {
		log.Printf("Impossible de récupérer le user par sont refreshToken %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the token is expires
	validToken, err := cfg.DB.GetRefreshTokenInfo(r.Context(), headerToken)
	if validToken.ExpiresAt.Before(time.Now().UTC()) {
		log.Printf("Token Expirer: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		log.Printf("Impossible de crée un nouveau Access Token: %s", err)
		return
	}

	utils.WriteJson(w, 200, NewToken{AccessToken: newAccessToken})
}

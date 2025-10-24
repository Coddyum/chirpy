package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type LoginUser struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type Refresh struct {
	ExpiresAt string `json:"expires_at"`
	RevokedAt string `json:"revoked_at"`
	Token     string `json:"token"`
}

func (cfg *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Decode et récupère les information de la request
	decoder := json.NewDecoder(r.Body)
	params := LoginUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// Récupère dans la base de donné le user grace a sont email
	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("DB error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// On vérifie que le password hash est correct
	checkPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !checkPassword {
		w.WriteHeader(401)
		return
	}

	// On crée un JWT d'acces valable 1h
	createJWT, err := auth.MakeJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		log.Printf("Error create jwt: %s", err)
		return
	}

	// On crée un refresh token valable 60Jour
	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error create refreshToken %s", err)
		return
	}

	// Ajoue des information du refreshToken a la db
	valid_util := time.Now().AddDate(0, 0, 60)
	refreshTokenDB, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		ExpiresAt: valid_util,
		Token:     newRefreshToken,
	})
	if err != nil {
		log.Printf("Error lors de l'ajoue des info de validité du refreshToken a la base de donner %s", err)
		return
	}

	// On écrit la réponse avec toute les informations sur le User
	utils.WriteJson(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        createJWT,
		RefreshToken: newRefreshToken,
	})

	// Information sur le refresh token
	utils.WriteJson(w, 200, Refresh{
		ExpiresAt: refreshTokenDB.ExpiresAt.String(),
		RevokedAt: refreshTokenDB.RevokedAt.Time.String(),
		Token:     refreshTokenDB.Token,
	})
}

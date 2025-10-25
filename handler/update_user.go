package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/Coddyum/chirpy/internal/utils"
)

type BodyParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *ApiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := BodyParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Le token est malformed or missing %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Le JWT n'est pas valid : %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.DB.GetUserFromAccessToken(r.Context(), userID)
	if err != nil {
		log.Printf("Impossible de récupérer le user par sont token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Erreur lors du hash du password %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	updatedUser, err := cfg.DB.UpdateUserInformation(r.Context(), database.UpdateUserInformationParams{
		ID:             userID,
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: newHashedPassword,
	})
	if err != nil {
		log.Printf("Erreur lors de l'update des information du user %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, 200, User{
		ID:           updatedUser.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
		Email:        updatedUser.Email,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
	})
}

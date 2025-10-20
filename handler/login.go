package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type LoginUser struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("DB error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	checkPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	//log.Printf("checkPassword: %v, err: %v", checkPassword, err)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var expiresIn int
	if params.ExpiresInSeconds <= 0 {
		expiresIn = 3600 // Max expiration time ( 1H )
	} else if params.ExpiresInSeconds > 3600 {
		expiresIn = 3600
	} else {
		expiresIn = params.ExpiresInSeconds
	}

	createJWT, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		log.Printf("Error create jwt: %s", err)
		return
	}

	if !checkPassword {
		w.WriteHeader(401)
		return
	}

	utils.WriteJson(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     createJWT,
	})

}

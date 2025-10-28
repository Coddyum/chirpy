package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/google/uuid"
)

type WebhookParams struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *ApiConfig) UpgradeUserWebHooks(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := WebhookParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	apiToken, err := auth.GetApiKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiToken != cfg.POLKA_KEY {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Erreur lors du parse du userID: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := cfg.DB.GetUserById(r.Context(), userID)
	if err != nil {
		log.Printf("User introuvable ou irrécupérable %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	UpdateErr := cfg.DB.UpgradesUser(r.Context(), database.UpgradesUserParams{
		ID:          user.ID,
		UpdatedAt:   time.Now(),
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	if UpdateErr != nil {
		log.Printf("Erreur lors de l'update %s", UpdateErr)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

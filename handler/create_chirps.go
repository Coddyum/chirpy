package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Coddyum/chirpy/internal/database"
	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type parameters struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type cleanedBody struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
}

func (cfg *ApiConfig) CreateChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		utils.WriteJson(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	data := cleanBodyString(params.Body)
	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   data,
		UserID: uuid.NullUUID{UUID: uuid.MustParse(params.UserId), Valid: true},
	})
	if err != nil {
		utils.WriteJson(w, 500, errorResponse{Error: "Failed to create chirp"})
		return
	}

	utils.WriteJson(w, 201, cleanedBody{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    params.UserId,
	})

}

func cleanBodyString(s string) string {
	splitedString := strings.Split(s, " ")
	cleanString := []string{}
	bannedWord := []string{"fornax", "sharbert", "kerfuffle"}

	for _, split := range splitedString {
		for _, ban := range bannedWord {
			lower := strings.ToLower(split)
			if strings.Contains(lower, ban) {
				split = strings.ReplaceAll(lower, ban, "****")
				break
			}
		}
		cleanString = append(cleanString, split)
	}

	joinString := strings.Join(cleanString, " ")
	return joinString
}

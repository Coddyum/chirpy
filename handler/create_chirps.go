package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type parameters struct {
	Body string `json:"body"`
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

	// Step 1: on récupère le bearer Token dans le Authorization du header.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("(create_chirps L.45) Impossible de récupérer le bearerToken : %s", err)
		return
	}

	// Step 2: On vérifie le que le token est bien valid
	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("(Create_Chirps L.52) Le token n'est pas valid %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Step 4: je commence a vérifié que les chirp sont valid avec la length
	if len(params.Body) > 140 {
		utils.WriteJson(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	// Step 5: je m'assure que le chirp ne contient pas de contenu interdu
	data := cleanBodyString(params.Body)

	// Step 6: je crée le chirp et je l'ajoute a la base de donner
	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   data,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		log.Printf("(Create_Chrips L.72) Impossible de crée le chirp %s", err)
		utils.WriteJson(w, 500, errorResponse{Error: "Failed to create chirp"})
		return
	}

	// Step 7: j'écrie la réponse a envoyer au user !
	utils.WriteJson(w, 201, cleanedBody{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    userID.String(),
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

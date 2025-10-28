package handler

import (
	"log"
	"net/http"

	"github.com/Coddyum/chirpy/internal/auth"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) DeleteChirpsHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Erreur lors du Parsing du chirp Id : %s", err)
		return
	}

	// Récupération du Access Token dans le header.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Impossible de récupérer le AccessToken (mal former ou manquant): %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Vérification du JWT et récupération du UserID lié a ce JWT Access Token
	validToken, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Erreur lors de la vérification du JWT Token %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Récupération du Chirp dans la req et vérifi qu'il existe
	chirp, err := cfg.DB.SelectOneChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Le chirp n'existe pas ou est impossible a récupérer %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Si l'author du chirp n'est pas celui qui en fait la demande return
	if chirp.UserID.UUID != validToken {
		log.Printf("L'author du chirp n'est pas le même que celui qui fait la request %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Delete chips
	deleteErr := cfg.DB.DeleteOneChirp(r.Context(), database.DeleteOneChirpParams{
		ID:     chirp.ID,
		UserID: uuid.NullUUID{UUID: validToken, Valid: true},
	})
	if deleteErr != nil {
		log.Printf("Error lors du delete du chirp %s", deleteErr)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

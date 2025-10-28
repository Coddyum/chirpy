package handler

import (
	"net/http"

	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerSelectOneChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DB.SelectOneChirp(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(404)
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	utils.WriteJson(w, 200, Chirps{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.UUID,
	})
}

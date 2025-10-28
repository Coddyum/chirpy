package handler

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Coddyum/chirpy/internal/database"
	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type Chirps struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) HandlerSelectAllChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	var chirp []database.Chirp
	var err error

	if id != "" {
		authorId, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}

		chirp, err = cfg.DB.SelectChirpByAuthor(r.Context(), uuid.NullUUID{UUID: authorId, Valid: true})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "Chirp not found", http.StatusNotFound)
			return
		}

	} else {
		chirp, err = cfg.DB.SelectAllChirps(r.Context())
		if err != nil {
			log.Fatalf("Erreur lors de la récupération de chirp %s", err)
			return
		}
	}

	if sorting == "desc" {
		sort.Slice(chirp, func(i, j int) bool {
			return chirp[i].CreatedAt.After(chirp[j].CreatedAt)
		})
	}

	var data []Chirps

	for _, c := range chirp {
		data = append(data,
			Chirps{
				Id:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserId:    c.UserID.UUID,
			})
	}
	utils.WriteJson(w, 200, data)
}

package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Coddyum/chirpy/internal/utils"
	"github.com/google/uuid"
)

type Chrips struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) HandlerSelectAllChirps(w http.ResponseWriter, r *http.Request) {
	chrip, err := cfg.DB.SelectAllChirps(r.Context())
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de chirp %s", err)
		return
	}

	var data []Chrips

	for _, c := range chrip {
		data = append(data,
			Chrips{
				Id:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserId:    c.UserID.UUID,
			})
	}
	utils.WriteJson(w, 200, data)
}

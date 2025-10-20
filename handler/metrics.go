package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Coddyum/chirpy/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	JWTSecret      string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Affiche le count du nombre de hits
func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, count)))
}

// Reset le compteur de visite
func (cfg *ApiConfig) ResetMetricHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// reset les metrics avant la suppression
	count := cfg.fileserverHits.Load()
	cfg.fileserverHits.CompareAndSwap(count, 0)

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to delete all users", http.StatusInternalServerError)
		return
	}

	// confirmation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted successfully"))
}

package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Incremente le hits x qu'un user var sur la main page (/app)
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Affiche le count du nombre de hits
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", count)))
}

// Reset le compteur de visite
func (cfg *apiConfig) resetMetricHandler(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	cfg.fileserverHits.CompareAndSwap(count, 0)
}

// Affiche le status de notre app
func readlinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))
	mux.Handle("/assets", fileServer)

	// Metric
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.resetMetricHandler)

	// Status off app health
	mux.HandleFunc("/healthz", readlinessHandler)

	srv.ListenAndServe()
}

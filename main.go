package main

import (
	"net/http"

	"github.com/Coddyum/chirpy/handler"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := &handler.ApiConfig{}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fileServer)))
	mux.Handle("/assets", fileServer)

	// Metric
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetMetricHandler)

	// Status off app health
	mux.HandleFunc("GET /api/healthz", handler.ReadlinessHandler)

	// Json
	mux.HandleFunc("POST /api/validate_chirp", handler.ValidatedChirp)

	srv.ListenAndServe()
}

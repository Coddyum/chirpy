package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// type returnValue struct {
// 	Valid bool `json:"valid"`
// }

type cleanedBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func writeJson(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

func ValidatedChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		writeJson(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	data, _ := cleanBodyString(params.Body)
	if len(data) > 1 {
		writeJson(w, 200, cleanedBody{CleanedBody: data})
		return
	}

	// writeJson(w, 200, returnValue{Valid: true})
}

func cleanBodyString(s string) (string, bool) {
	splitedString := strings.Split(s, " ")
	cleanString := []string{}
	bannedWord := []string{"fornax", "sharbert", "kerfuffle"}
	replace := false

	for _, split := range splitedString {
		for _, ban := range bannedWord {
			lower := strings.ToLower(split)
			if strings.Contains(lower, ban) {
				split = strings.ReplaceAll(lower, ban, "****")
				replace = true
				break
			}
		}
		cleanString = append(cleanString, split)
	}

	joinString := strings.Join(cleanString, " ")
	return joinString, replace
}

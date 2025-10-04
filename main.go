package main

import "net/http"

type Server struct {
	Addr    string
	Handler http.ServeMux
}

func main() {
	mux := http.NewServeMux()

	http.ListenAndServe(":8080", mux)
}

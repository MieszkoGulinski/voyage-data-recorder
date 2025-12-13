package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func startHTTPServer() {
	r := chi.NewRouter()

	// TODO complete
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	http.ListenAndServe(":8080", r)
}
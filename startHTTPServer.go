package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func startHTTPServer(db *gorm.DB) {
	r := chi.NewRouter()

	// TODO complete
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	http.ListenAndServe(":8080", r)
}
package viewerhttp

import (
	"datalogger/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func generateApiListHandler[T database.WithTimestamp](db *gorm.DB,
	model *T,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var lastTimestamp int64 = 0
		if v := r.URL.Query().Get("lastTimestamp"); v != "" {
			if parsedTimestamp, err := strconv.ParseInt(v, 10, 64); err == nil {
				lastTimestamp = parsedTimestamp
			}
		}

		rows, _, _, err := database.QueryWithPagination(db, model, lastTimestamp, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func StartHTTPServer(db *gorm.DB) {
	r := chi.NewRouter()

	// TODO write HTML renderer
	r.Get("/api/weather", generateApiListHandler(db, &database.Weather{}))
	r.Get("/api/positions", generateApiListHandler(db, &database.Position{}))
	r.Get("/api/battery",  generateApiListHandler(db, &database.Battery{}))

	http.ListenAndServe(":8080", r)
}
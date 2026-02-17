package listeners

import (
	"context"
	"datalogger/writer"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type MessageRequest struct {
	Message string `json:"message"`
}

func StartHTTPListener(ctx context.Context, port int, diagnostics bool, channelsSet writer.ChannelsSet) error {
	r := chi.NewRouter()

	// API route handlers
	r.Post("/navtex", func(w http.ResponseWriter, r *http.Request) {
		var body MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if body.Message == "" {
			http.Error(w, "Message cannot be empty", http.StatusBadRequest)
			return
		}
		channelsSet.NavtexCh <- body.Message
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/position", func(w http.ResponseWriter, r *http.Request) {
		var body writer.PositionMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if body.Latitude > 90 || body.Latitude < -90 || body.Longitude > 180 || body.Longitude < -180 {
			http.Error(w, "Invalid position", http.StatusBadRequest)
			return
		}
		channelsSet.PositionCh <- body
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/text", func(w http.ResponseWriter, r *http.Request) {
		var body MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if body.Message == "" {
			http.Error(w, "Message cannot be empty", http.StatusBadRequest)
			return
		}
		channelsSet.TextNoteCh <- body.Message
		w.WriteHeader(http.StatusOK)
	})

	addr := fmt.Sprintf(":%d", port)

	// Start server with graceful shutdown using context
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if diagnostics {
		fmt.Printf("HTTP listener starting on port %d\n", port)
	}

	// Channel to capture server errors
	errCh := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		return server.Shutdown(shutdownCtx)

	case err := <-errCh:
		return err
	}

}

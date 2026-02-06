package listeners

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func StartHTTPListener(ctx context.Context, port int, diagnostics bool) error {
	r := chi.NewRouter()

	// Attach API route handlers
	r.Get("/navtex", func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
	r.Get("/position", func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
	r.Get("/text", func(w http.ResponseWriter, r *http.Request) {
		// TODO
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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		return server.Shutdown(shutdownCtx)

	case err := <-errCh:
		return err
	}

}

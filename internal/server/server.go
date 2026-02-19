package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devaloi/htmxapp/internal/handler"
	"github.com/devaloi/htmxapp/internal/store"
	"github.com/devaloi/htmxapp/internal/tmpl"
)

// Run starts the HTTP server with graceful shutdown.
func Run(cfg Config) error {
	renderer, err := tmpl.New()
	if err != nil {
		return fmt.Errorf("initializing templates: %w", err)
	}

	memStore := store.NewMemory()
	if cfg.Seed {
		memStore.Seed()
		slog.Info("seeded sample contacts")
	}

	h := handler.New(memStore, renderer)
	routes := h.Routes()

	// Apply middleware stack: RequestID → Recovery → Logging → routes
	stack := handler.RequestIDMiddleware(handler.Recovery(handler.Logging(routes)))

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      stack,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", cfg.Addr())
		if serveErr := srv.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
		close(errCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		slog.Info("shutting down", "signal", sig.String())
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	slog.Info("server stopped")
	return nil
}

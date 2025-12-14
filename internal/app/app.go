package app

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{
		httpServer: &http.Server{},
	}
}

func (a *App) Start(ctx context.Context) error {
	fmt.Println("Starting application...")

	a.httpServer.Addr = ":8080"

	return runHTTPServer(ctx, a.httpServer)
}

func (a *App) Stop(ctx context.Context) error {
	return shutdownHTTPServer(ctx, a.httpServer)
}

func runHTTPServer(ctx context.Context, server *http.Server) error {
	fmt.Println("Starting HTTP server on", server.Addr)

	errChan := make(chan error, 1)

	go func() {
		fmt.Println("HTTP server is running...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return shutdownHTTPServer(ctx, server)
	case err := <-errChan:
		return err
	}
}

func shutdownHTTPServer(ctx context.Context, server *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}

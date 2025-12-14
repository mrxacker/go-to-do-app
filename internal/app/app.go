package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mrxacker/go-to-do-app/internal/config"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config

	httpServer *http.Server
	grpcServer *grpc.Server
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	grpcSrv := grpc.NewServer()

	httpSrv := &http.Server{
		Addr: cfg.HTTPAddr,
	}

	return &App{
		cfg:        cfg,
		httpServer: httpSrv,
		grpcServer: grpcSrv,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	fmt.Println("Starting application...")

	a.httpServer.Addr = ":8080"

	errChan := make(chan error, 2)

	go func() {
		err := runGRPCServer(ctx, a.grpcServer, a.cfg.GRPCAddr)
		errChan <- err
	}()

	go func() {
		err := runHTTPServer(ctx, a.httpServer)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return a.Stop(ctx)
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("Stopping application...")
	err := shutdownHTTPServer(ctx, a.httpServer)
	if err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	a.grpcServer.GracefulStop()
	return nil
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
		fmt.Println("Stopping http server")
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

func runGRPCServer(ctx context.Context, grpcSrv *grpc.Server, grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcAddr, err)
	}

	fmt.Println("Starting gRPC server on", lis.Addr().String())

	errChan := make(chan error, 1)

	go func() {
		fmt.Println("gRPC server is running...")
		if err := grpcSrv.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Stopping gRPC server")
		grpcSrv.GracefulStop()
		return nil
	case err := <-errChan:
		return err
	}
}

package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	internal_http "github.com/mrxacker/go-to-do-app/internal/adapters/http"
	"github.com/mrxacker/go-to-do-app/internal/config"
	"github.com/mrxacker/go-to-do-app/internal/infrastructure/postgres"
	"github.com/mrxacker/go-to-do-app/internal/logger"
	"github.com/mrxacker/go-to-do-app/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	cfg    *config.Config
	logger *zap.Logger

	todoUC *usecase.TodoUsecase

	httpServer *http.Server
	grpcServer *grpc.Server
	db         *sql.DB

	wg sync.WaitGroup
}

func NewApp() (*App, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	l, err := logger.NewLogger(cfg.ENV)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database connection
	db, err := postgres.NewPostgresDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repositories and use cases
	todoRepo := postgres.NewTodoRepo(db)
	todoUC := usecase.NewTodoUsecase(todoRepo)
	userRepo := postgres.NewUserRepo(db)
	userUC := usecase.NewUserUseCase(userRepo)

	// Initialize HTTP handlers
	httpRouter := initHandlers(todoUC, userUC)

	// Initialize servers
	grpcSrv := grpc.NewServer()
	httpSrv := &http.Server{Addr: ":" + cfg.HTTPAddr, Handler: httpRouter}

	// Return the application instance
	return &App{
		cfg:        cfg,
		todoUC:     todoUC,
		httpServer: httpSrv,
		grpcServer: grpcSrv,
		logger:     l.Logger,
		db:         db,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("Starting application",
		zap.String("http_addr", a.cfg.HTTPAddr),
		zap.String("grpc_addr", a.cfg.GRPCAddr),
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		errCh <- a.runHTTP(ctx)
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		errCh <- a.runGRPC(ctx)
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			a.logger.Error("application shutdown", zap.Error(err))
			cancel()
			_ = a.shutdown()
			return err
		}
	}

	return nil
}

func (a *App) shutdown() error {
	a.logger.Info("Shutting down application")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	_ = a.shutdownHTTP()
	a.shutdownGRPC()

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-shutdownCtx.Done():
		return shutdownCtx.Err()
	}

	return nil
}

func (a *App) runHTTP(ctx context.Context) error {
	if a.httpServer.Addr == "" {
		return errors.New("http address is empty")
	}

	a.logger.Info("Starting HTTP server", zap.String("http_addr", a.httpServer.Addr))

	go func() {
		<-ctx.Done()
		a.logger.Info("Shutting down HTTP server", zap.String("http_addr", a.httpServer.Addr))
		_ = a.shutdownHTTP()
	}()

	if err := a.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return context.Canceled
		}
		a.logger.Error("application shutdown", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) shutdownHTTP() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}

func (a *App) runGRPC(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+a.cfg.GRPCAddr)
	if err != nil {
		a.logger.Error("application shutdown", zap.Error(err))
		return err
	}

	a.logger.Info("Starting gRPC server", zap.String("grpc_addr", a.cfg.GRPCAddr))

	go func() {
		<-ctx.Done()
		a.logger.Info("Shutting down gRPC server", zap.String("grpc_addr", a.cfg.GRPCAddr))
		a.shutdownGRPC()
	}()

	if err := a.grpcServer.Serve(lis); err != nil {
		a.logger.Error("application shutdown", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) shutdownGRPC() {
	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		a.grpcServer.Stop()
	}
}

func initHandlers(todoUC *usecase.TodoUsecase, userUC *usecase.UserUseCase) *gin.Engine {
	todoHandler := internal_http.NewTodoHandler(todoUC)
	userHandler := internal_http.NewUserHandler(userUC)
	r := gin.Default()
	api := r.Group("/api/v1")
	todoHandler.RegisterRoutes(api.Group("/todos"))
	userHandler.RegisterRoutes(api.Group("/users"))
	return r
}

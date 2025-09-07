package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"mediary/cmd/internal"
	"mediary/config"
	internaldomain "mediary/internal"
	"mediary/internal/httphandler"
	"mediary/internal/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	var envFile string

	flag.StringVar(&envFile, "env", "", "Environment Variables filepath")
	flag.Parse()

	errC, err := run(envFile)
	if err != nil {
		log.Fatalf("Couldn't run: %v", err)
	}

	if runErr := <-errC; runErr != nil {
		if e, ok := runErr.(*internaldomain.Error); ok {
			logger, _ := zap.NewProduction() // fallback logger in case logger not available
			defer logger.Sync()
			logger.Error("application error", e.Fields()...)
		} else {
			log.Printf("unexpected error: %v", runErr)
		}
		os.Exit(1)
	}
}

func run(envFile string) (<-chan error, error) {
	cfg, err := config.Load(envFile)
	if err != nil {
		return nil, internaldomain.NewConfigError("Failed to load config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, internaldomain.NewConfigError("failed to initialize logger: %v", err)
	}

	if err := internal.MigratePostgresql(cfg, "db/migrations"); err != nil {
		return nil, internaldomain.WrapDatabaseError(err, "internal.MigratePostgresql failed")
	}

	pool, err := internal.NewPostgreSQL(cfg)
	if err != nil {
		return nil, internaldomain.WrapDatabaseError(err, "internal.NewPostgreSQL failed")
	}
	if pool != nil {

	}

	srv := internal.NewServer(
		cfg,
		[]httphandler.Handler{},

		chimiddleware.RealIP,
		middleware.TracingMiddleware,
		middleware.LoggingMiddleware(logger),
		middleware.RecoveryMiddleware(logger),
		middleware.AuthMiddleware(cfg),
		middleware.CORSMiddleware(),
		middleware.MetricsMiddleware(),
		chimiddleware.Timeout(15*time.Second),
	)

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Shutdown
	go func() {
		<-ctx.Done()
		logger.Info("Shutdown signal recieved")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			_ = logger.Sync()

			pool.Close()
			srv.Close()

			stop()
			cancel()
		}()
		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- internaldomain.WrapServerError(err, "server shutdown failed")
		}

		logger.Info("Shutdown completed")
	}()

	// Start
	go func() {
		logger.Info("Listening and serving", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- internaldomain.WrapServerError(err, "ListenAndServe failed")
		}
	}()

	return errC, nil
}

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/pintoter/persons/docs"
	"github.com/pintoter/persons/pkg/database/postgres"
	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/query/internal/config"
	dbrepo "github.com/pintoter/persons/services/query/internal/repository/db"
	"github.com/pintoter/persons/services/query/internal/server"
	"github.com/pintoter/persons/services/query/internal/service"
	"github.com/pintoter/persons/services/query/internal/transport"
)

// @title           			Persons
// @version         			1.0
// @description     			REST API service Persons

// @contact.name   				Vlad Yurasov
// @contact.email  				meine23@yandex.ru

// @host      					persons-query:8080
// @BasePath  					/api/v1

func Run() {
	ctx := context.Background()

	cfg := config.Get()

	syncLogger := initLogger(ctx, cfg)
	defer syncLogger()

	db, err := postgres.New(&cfg.DB)
	if err != nil {
		logger.FatalKV(ctx, "Failed connect database", "err", err)
	}

	repo := dbrepo.New(db)

	service := service.New(repo)
	handler := transport.NewHandler(service)
	server := server.New(handler, &cfg.HTTP)

	server.Run()
	logger.InfoKV(ctx, "Starting server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	select {
	case <-quit:
		logger.InfoKV(ctx, "Starting gracefully shutdown")
	case err = <-server.Notify():
		logger.FatalKV(ctx, "Failed starting server", "err", err.Error())
	}

	if err := server.Shutdown(); err != nil {
		logger.FatalKV(ctx, "Failed shutdown server", "err", err.Error())
	}
}

func initLogger(ctx context.Context, cfg *config.Config) (syncFn func()) {
	loggingLevel := zap.InfoLevel
	if cfg.Project.Level == logger.DebugLevel {
		loggingLevel = zap.DebugLevel
	}

	loggerConfig := zap.NewProductionEncoderConfig()

	loggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(loggerConfig),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	notSuggaredLogger := zap.New(consoleCore)

	sugarLogger := notSuggaredLogger.Sugar()
	logger.SetLogger(sugarLogger.With(
		"service", cfg.Project.Name,
	))

	return func() {
		_ = notSuggaredLogger.Sync()
	}
}

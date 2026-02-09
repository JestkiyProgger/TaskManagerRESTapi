package main

import (
	"ProjectManagementAPI/internal/config"
	tcrud "ProjectManagementAPI/internal/http-server/handlers/task/crud"
	ucrud "ProjectManagementAPI/internal/http-server/handlers/user/crud"
	mwLogger "ProjectManagementAPI/internal/http-server/middleware/logger"
	"ProjectManagementAPI/internal/lib/logger/handlers/slogpretty"
	"ProjectManagementAPI/internal/lib/logger/sl"
	"ProjectManagementAPI/internal/storage/postgre"
	tasksvc "ProjectManagementAPI/internal/task"
	usersvc "ProjectManagementAPI/internal/user"
	"log/slog"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv

	cfg := config.MustLoad()

	// init logger: slog

	logger := setupLogger(cfg.Env)

	logger.Info("starting task-management", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// init storage: PostgreSQL

	storage, err := postgre.New(cfg.Postgres.DSN)
	if err != nil {
		logger.Error("failed to initialize postgresql storage", sl.Err(err))
		os.Exit(1)
	}

	defer func(storage *postgre.Storage) {
		err := storage.Close()
		if err != nil {
			logger.Error("failed to close postgresql storage", sl.Err(err))
			os.Exit(1)
		}
	}(storage)

	// init router: chi, "chi render"

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	userRepo := usersvc.NewRepository(storage.Db)
	userService := usersvc.NewService(userRepo)

	taskRepo := tasksvc.NewRepository(storage.Db)
	taskService := tasksvc.NewService(taskRepo)

	router.Post("/user", ucrud.CreateNew(logger, userService))
	router.Delete("/user/{id}", ucrud.NewDeleteHandler(logger, userService))

	router.Post("/task", tcrud.CreateNew(logger, taskService))
	router.Delete("/task/{id}", tcrud.NewDeleteHandler(logger, taskService))

	logger.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server stopped")
	// TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = setupPrettySlog()
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	color.NoColor = false

	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

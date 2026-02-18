package main

import (
	"ProjectManagementAPI/internal/config"
	taskHttp "ProjectManagementAPI/internal/http-server/handlers/task"
	userHttp "ProjectManagementAPI/internal/http-server/handlers/user"
	mwLogger "ProjectManagementAPI/internal/http-server/middleware/logger"
	"ProjectManagementAPI/internal/lib/logger/sl"
	taskRepository "ProjectManagementAPI/internal/repository/postgres/task"
	userRepository "ProjectManagementAPI/internal/repository/postgres/user"
	"ProjectManagementAPI/internal/storage/postgre"
	taskService "ProjectManagementAPI/internal/usecase/task"
	userService "ProjectManagementAPI/internal/usecase/user"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting task-management", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

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

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	userRepo := userRepository.NewUserRepository(storage.Db)
	taskRepo := taskRepository.NewTaskRepository(storage.Db)

	userServ := userService.NewUserService(userRepo)
	taskServ := taskService.NewTaskService(taskRepo)

	userHandler := userHttp.NewHandler(logger, userServ)
	taskHandler := taskHttp.NewHandler(logger, taskServ)

	router.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.Create)
		r.Delete("/{id}", taskHandler.Delete)
		r.Get("/{id}", taskHandler.GetByID)
	})

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Delete("/{id}", userHandler.Delete)
		r.Get("/{id}", userHandler.GetByID)
	})

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
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

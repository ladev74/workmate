package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"link-service/internal/handler"
	"link-service/internal/logger"
	"link-service/internal/repository"
	"link-service/internal/service"
)

type Config struct {
	Host            string        `env:"HTTP_HOST" env-required:"true"`
	Port            int           `env:"HTTP_PORT" env-required:"true"`
	Timeout         time.Duration `env:"HTTP_OPERATION_TIMEOUT" env-required:"true"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-required:"true"`
}

func NewRouter(cfgLogger *logger.Config, cfgServer *Config, log *zap.Logger, repo repository.Repository) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.MiddlewareLogger(log, cfgLogger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: set a timeout in the config
	srv := service.New(repo, 3*time.Second, log)

	router.Post("/links", handler.ProcessLinks(srv, cfgServer.Timeout, log))

	return router
}

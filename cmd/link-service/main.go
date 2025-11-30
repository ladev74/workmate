package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"link-service/internal/config"
	"link-service/internal/logger"
	filesystem "link-service/internal/repository/file_system"
	"link-service/internal/server"
)

// TODO: rename repo
// TODO: create logger interface
// TODO: flag -race

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cfgPath := fetchConfigPath()
	if cfgPath == "" {
		stdlog.Fatal("config file path must be specified")
	}

	cfg, err := config.New(cfgPath)
	if err != nil {
		stdlog.Fatalf("cannot initialize config: %v", err)
	}

	log, err := logger.New(&cfg.Logger)
	if err != nil {
		stdlog.Fatalf("cannot initialize logger: %v", err)
	}
	defer log.Sync()

	storage := filesystem.New(log)
	// TODO: вынести в конфиг
	err = storage.Init("./data", "data.json", "temp.json")
	if err != nil {
		log.Fatal("cannot initialize storage: %v", zap.Error(err))
		return
	}

	srv := server.New(ctx, &cfg.Logger, &cfg.HTTPServer, log, storage)

	go func() {
		log.Info("starting http server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	//shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	//defer shutdownCancel()

	log.Info("received shutdown signal")

	//err = srv.Shutdown(shutdownCtx)
	//if err != nil {
	//	log.Error("failed to shutdown server", zap.Error(err))
	//}

	log.Info("application shutdown completed successfully")
}

func fetchConfigPath() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config_path", "", "path to config file")
	flag.Parse()

	return cfgPath
}

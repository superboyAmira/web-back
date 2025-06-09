package main

import (
	"backend/internal/api"
	repo "backend/internal/repository"
	"backend/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := repo.NewInMemoryRepo()
	svc := service.NewTournamentService(repo)
	h := api.NewHandler(svc, logger)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}

	go func() {
		logger.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"proxy/internal/config"
	"proxy/internal/delivery/handlers"
	"proxy/internal/repository"
	"proxy/internal/server"
	"proxy/internal/usecase"
	logger "proxy/pkg/logger"
	"syscall"
)

func main() {

	cfg := config.Load()

	logger := logger.NewLogger()

	white, _ := repository.NewIpRepoImpl("white")
	black, _ := repository.NewIpRepoImpl("black")
	gray, _ := repository.NewIpRepoImpl("gray")
	rlStore, _ := repository.NewRateLimitStore("configs/ratelimit.json")

	whiteUC := usecase.NewHTTPUseCase(white)
	blackUC := usecase.NewHTTPUseCase(black)
	grayUC := usecase.NewHTTPUseCase(gray)
	rlUC := usecase.NewRateLimitService(rlStore, cfg)

	whiteHandler := handlers.NewHandler(whiteUC, logger)
	blackHandler := handlers.NewHandler(blackUC, logger)
	grayHandler := handlers.NewHandler(grayUC, logger)

	rlHandler := handlers.NewRateLimitHandler(rlUC)

	checkHandler := handlers.NewCheckIpHandler(whiteUC, blackUC, grayUC, logger)

	srv := server.NewServer(cfg, logger, whiteHandler, grayHandler, blackHandler, checkHandler, rlHandler, rlUC)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}

}

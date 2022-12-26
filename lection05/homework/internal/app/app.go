package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cr00z/chat/internal/controller/http"
	"github.com/cr00z/chat/internal/infrastructure/repository/memory"
	"github.com/cr00z/chat/internal/service"
)

func Run() {
	repos := repository.New()
	services := service.New(repos)
	handlers := handler.New(services)

	srv := new(Server)
	go func() {
		if err := srv.Run("5000", handlers.InitRoutes()); err != nil {
			// TODO:
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// TODO:
	log.Println("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// TODO:
	log.Println("signal received:", <-quit, ", server shutting down")

	if err := srv.Shutdown(context.TODO()); err != nil {
		// TODO:
		log.Printf("error occured on server shutting down: %s", err.Error())
	}
}

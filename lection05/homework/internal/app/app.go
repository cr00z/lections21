package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cr00z/goSimpleChat/internal/controller/handler"
	repository "github.com/cr00z/goSimpleChat/internal/infrastructure/repository/memory"
	"github.com/cr00z/goSimpleChat/internal/service"
)

func Run() {
	repos, err := repository.New()
	if err != nil {
		// TODO: change logger
		log.Fatalf("error occured while init repository: %s", err.Error())
	}

	services := service.New(repos)
	handlers := handler.New(services)

	srv := new(Server)
	go func() {
		if err := srv.Run("5000", handlers.InitRoutes()); err != nil {
			// TODO: change logger
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// TODO: change logger
	log.Println("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// TODO: change logger
	log.Println("signal received:", <-quit, ", server shutting down")

	if err := srv.Shutdown(context.TODO()); err != nil {
		// TODO: change logger
		log.Printf("error occured on server shutting down: %s", err.Error())
	}
}

package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ozon_link_shortener/internal/service"

	"github.com/joho/godotenv"
)

func InitEnv() *service.Config {
	cfg := service.Config{}
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Error loading env file %s: %v", envFile, err)
		}
	} else {
		candidateFiles := []string{".env", ".env"}
		loaded := false
		for _, file := range candidateFiles {
			if err := godotenv.Load(file); err == nil {
				loaded = true
				break
			}
		}
		if !loaded {
			log.Printf("No env file loaded (tried %v); using process environment", candidateFiles)
		}
	}
	cfg.Storage = os.Getenv("STORAGE")
	if cfg.Storage == "" {
		log.Println("Using cache storage")
	}
	cfg.DBAddress = os.Getenv("DATABASE_ADDRESS")
	if cfg.DBAddress == "" {
		log.Println("No address")
	}
	return &cfg
}

func GracefulShutdown(s *service.Service, server *http.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		fmt.Println("\nShutdown...")
		cancel()
		if s.Conf.Storage == "postgres" {
			err := s.DB.Close()
			if err != nil {
				log.Println(err)
			}
		}
		server.Shutdown(ctx)
	}()
}

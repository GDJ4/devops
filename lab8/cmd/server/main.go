package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todoapp/internal/config"
	"todoapp/internal/todo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := mongo.Connect(rootCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("connect to mongo: %v", err)
	}

	if err := client.Ping(rootCtx, nil); err != nil {
		log.Fatalf("ping mongo: %v", err)
	}
	log.Printf("connected to MongoDB at %s", cfg.MongoURI)

	repo := todo.NewRepository(client, cfg.MongoDB)
	server := todo.NewServer(repo)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: server.Router(),
	}

	go func() {
		<-rootCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Println("shutting down http server...")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
		if err := client.Disconnect(shutdownCtx); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()

	log.Printf("listening on :%s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}

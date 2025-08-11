package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"wb-techschool/L0/httpserver"
	"wb-techschool/L0/internal/kafka"
	"wb-techschool/L0/internal/repo/postgres"
	"wb-techschool/L0/internal/usecases"

	"github.com/go-chi/chi/v5"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	db, err := postgres.NewDB(postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	defer db.Close()

	repo := postgres.NewOrderRepo(db)
	cache := usecases.NewOrderCache()
	if err := cache.Warmup(context.Background(), repo); err != nil {
		log.Printf("cache warmup error: %v", err)
	}

	brokers := getEnv("KAFKA_BROKERS", "localhost:29092")
	topic := getEnv("KAFKA_TOPIC", "orders")
	group := getEnv("KAFKA_GROUP_ID", "orders-consumer")
	kcfg := kafka.Config{
		Brokers:  []string{brokers},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6,
	}
	consumer := kafka.NewConsumer(kcfg, repo, cache)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("kafka consumer stopped: %v", err)
		}
	}()

	srv := &httpserver.Server{Repo: repo, Cache: cache}
	var r chi.Router = httpserver.NewRouter(srv)

	addr := getEnv("HTTP_ADDR", ":8081")
	log.Printf("starting HTTP server on %s", addr)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httpserver.CreateAndRunServer(r, addr)
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutting down...")
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}

	time.Sleep(500 * time.Millisecond)
}

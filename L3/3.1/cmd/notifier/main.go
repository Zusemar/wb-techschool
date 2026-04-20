package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"notifier/internal/handlers"
	"notifier/internal/queue"
	"notifier/internal/repo"
	"notifier/internal/sender"
	"notifier/internal/usecases"
	"notifier/internal/worker"
)

func main() {
	rabbitURL := envOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	serverAddr := envOrDefault("SERVER_ADDR", ":8080")

	// RabbitMQ client
	rabbitClient, err := rabbitmq.NewClient(rabbitmq.ClientConfig{
		URL:            rabbitURL,
		ConnectionName: "notifier",
		ConnectTimeout: 10 * time.Second,
		Heartbeat:      10 * time.Second,
		ReconnectStrat: retry.Strategy{Attempts: 10, Delay: 1 * time.Second, Backoff: 2.0},
		ProducingStrat: retry.Strategy{Attempts: 3, Delay: 500 * time.Millisecond, Backoff: 2.0},
		ConsumingStrat: retry.Strategy{Attempts: 1, Delay: 0, Backoff: 1.0},
	})
	if err != nil {
		log.Fatalf("connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	// Declare exchange
	if err := rabbitClient.DeclareExchange(
		queue.ExchangeName, "direct", true, false, false, nil,
	); err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	// Declare queue and bind to exchange
	if err := rabbitClient.DeclareQueue(
		queue.QueueName, queue.ExchangeName, queue.RoutingKey,
		true, false, true, nil,
	); err != nil {
		log.Fatalf("declare queue: %v", err)
	}

	// Layers
	repository := repo.NewMemoryRepo()
	multiSender := sender.NewMultiSender()
	notifier := usecases.NewNotifier(repository)
	publisher := queue.NewPublisher(rabbitClient)
	scheduler := worker.NewScheduler(repository, publisher)
	consumer := queue.NewConsumer(rabbitClient, repository, multiSender)

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start background scheduler
	go scheduler.Run(ctx)

	// Start RabbitMQ consumer
	go func() {
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	// HTTP server
	engine := ginext.New("debug")
	handler := handlers.NewHandler(notifier)
	handler.RegisterRoutes(engine)

	log.Printf("server listening on %s", serverAddr)
	if err := engine.Run(serverAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

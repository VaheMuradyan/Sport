package main

import (
	"context"
	"github.com/VaheMuradyan/Sport/generator"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("🚀 Starting Odds Generator...")

	// Create odds generator
	gen, err := generator.NewCoefficientGenerator("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create odds generator: %v", err)
	}
	defer gen.Stop()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start generator in goroutine
	go gen.Start(ctx)

	// Wait for shutdown signal
	<-sigChan
	log.Println("🛑 Shutdown signal received, stopping generator...")
	cancel()
}

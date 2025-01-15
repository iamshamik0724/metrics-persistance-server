package main

import (
	"context"
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/boot"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	config, configErr := config.LoadConfig()
	if configErr != nil {
		panic("Failed to load config")
	}

	stopChannel := make(chan struct{})

	// Initialize Service
	udpErr := boot.Initialize(config, ctx, stopChannel)
	if udpErr != nil {
		log.Fatalf("Failed to Initialize Service: %v\n", udpErr)
	}

	go handleGracefulShutdown(cancel, stopChannel)

	log.Println("Service Initialized Successfully")

	<-ctx.Done()
	log.Println("Application shutting down.")
}

func handleGracefulShutdown(cancel context.CancelFunc, stopChannel chan struct{}) {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	log.Println("Shutdown signal received.")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(stopChannel)
	}()
	wg.Wait()
	bufferDuration := 5 * time.Second // Add a 5-second buffer time before cancelling context
	time.Sleep(bufferDuration)
	cancel()
}

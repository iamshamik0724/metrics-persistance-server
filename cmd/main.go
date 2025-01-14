package main

import (
	"context"
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/boot"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	config, configErr := config.LoadConfig()
	if configErr != nil {
		panic("Failed to load config")
	}

	go handleGracefulShutdown(cancel)

	// Initialize UDP connection
	udpErr := boot.InitializeUdpConnection(config, ctx)
	if udpErr != nil {
		log.Fatalf("Failed to create UDP connection: %v\n", udpErr)
	}

	log.Println("UDP connection successfully established.")

	<-ctx.Done()
	log.Println("Application shutting down.")
}

func handleGracefulShutdown(cancel context.CancelFunc) {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	log.Println("Shutdown signal received.")
	cancel()
}

//TODO: handle connection close on server shutdown

//Initialize DB

//Initilize repo

//Initialize Services

//Initialize handler

// bring up the server

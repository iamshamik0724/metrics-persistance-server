package main

import (
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/boot"
)

func main() {
	config, configErr := config.LoadConfig()
	if configErr != nil {
		panic("Failed to load config")
	}

	udpError := boot.InitlizeUdpConnection(config)
	if udpError != nil {
		panic("Failed to create UDP connection")
	} else {
		log.Println("UDP successfull")
	}

	//Initialize DB

	//Initilize repo

	//Initialize Services

	//Initialize handler

	//Initialize UDP connection

	// Send the connection Init message

	// bring up the server

}

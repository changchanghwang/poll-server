package main

import (
	"log"

	"poll.ant/cmd/di"
	"poll.ant/internal/config"
)

func main() {
	server, err := di.InitializeServer()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := server.Run(":" + config.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

package main

import (
	"os"
)

func main() {
	migrateDatabase()

	port := os.Getenv("WEBSOCKET_PORT")
	if port == "" {
		port = "1887"
	}

	runHttpServer(port)
}

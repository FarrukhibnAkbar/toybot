package main

import (
	"log"

	"toybot/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("❌ App initialization failed: %v", err)
	}

	application.Run()
}

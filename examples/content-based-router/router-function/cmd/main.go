package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/holgeradam/eip/examples/content-based-router/router-function"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.RegisterEventFunctionContext(context.Background(), "/", router.Route); err != nil {
		log.Fatalf("funcframework.RegisterEventFunctionContext: %v\n", err)
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

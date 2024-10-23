// Package main is used to create a localhost server.
package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	_ "github.com/kimikimi714/jimiko"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

// see: https://cloud.google.com/functions/docs/create-deploy-http-go?hl=ja
func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

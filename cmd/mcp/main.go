package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "client-runaway-zenoti/internal/db" // Initialize DB connection
	"client-runaway-zenoti/internal/mcp"
)

func main() {
	// Get address from environment or use default
	addr := os.Getenv("MCP_ADDR")
	if addr == "" {
		addr = ":8090"
	}

	server := mcp.NewMCPServer(addr)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down MCP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("Starting MCP server on %s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

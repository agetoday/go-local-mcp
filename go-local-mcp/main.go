package main

import (
	"gocode/go-local-mcp/mcpRouter"
	"log"
)

func main() {
	s := mcpRouter.NewMCPServer()
	sseServer := s.ServeSSE("127.0.0.1:8111")
	log.Printf("SSE server listening on :8111")
	if err := sseServer.Start(":8111"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

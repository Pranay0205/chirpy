package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	const port = "8080"

	mux := http.NewServeMux()

	server := http.Server{Addr: ":" + port, Handler: mux}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🌟 HTTP Server Starting...")
	fmt.Println("📍 Address: http://localhost:8080")
	fmt.Println("⏰ Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("💡 Tip: Press Ctrl+C to stop")
	fmt.Println(strings.Repeat("=", 50))
	
	log.Printf("✅ Server is listening on port 8080...")
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
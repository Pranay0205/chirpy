package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	server := http.Server{Addr: ":" + port, Handler: mux}
	fmt.Println("Server Starting...")
	fmt.Println("Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Tip: Press Ctrl+C to stop")
	
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	log.Printf("Serving files from %s on port: %s...\n", filepathRoot, port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}



}
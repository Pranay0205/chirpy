package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)



func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") 

	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
func main() {

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	fmt.Println("🚀  Server Starting...")
	fmt.Println("🕒  Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("💡  Tip: Press Ctrl+C to stop")
	
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", healthHandler)

	server := &http.Server{Addr: ":" + port, Handler: mux}

	fmt.Printf("⚡  Serving files from %s on port: %s...\n", filepathRoot, port)

	log.Fatal(server.ListenAndServe()) 
}
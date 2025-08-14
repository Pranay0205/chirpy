package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		log.Printf("📃	%s %s - User-Agent: %s", r.Method, r.URL.String(),	r.UserAgent())
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	count := cfg.fileserverHits.Load()

	message := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, count)

	w.Write([]byte(message))
}


func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") 

	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func main() {

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	apiConfig := apiConfig{fileserverHits: atomic.Int32{},}

	fmt.Println("🚀  Server Starting...")
	fmt.Println("🕒  Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("💡  Tip: Press Ctrl+C to stop")
	
	mux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	
	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleResetMetrics)

	server := &http.Server{Addr: ":" + port, Handler: mux}

	fmt.Printf("⚡  Serving files from %s on port: %s...\n", filepathRoot, port)

	log.Fatal(server.ListenAndServe()) 
}
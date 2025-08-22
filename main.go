package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📃	%s %s - User-Agent: %s", r.Method, r.URL.String(), r.UserAgent())
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(200)

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

	w.WriteHeader(200)

	w.Write([]byte(http.StatusText(200)))
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("❌	failed to establish connection with db")
	}
	fmt.Println("💾  Database connection established")

	mux := http.NewServeMux()
	apiConfig := apiConfig{fileserverHits: atomic.Int32{}, dbQueries: database.New(db), platform: os.Getenv("PLATFORM"), secret: os.Getenv("SECRET"), polkaKey: os.Getenv("POLKA_KEY")}

	fmt.Println("🚀  Server Starting...")
	fmt.Println("🕒  Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("💡  Tip: Press Ctrl+C to stop")

	mux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleResetMetrics)

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("GET /api/chirps", apiConfig.handleChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handleChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handleDeleteChirps)
	mux.HandleFunc("POST /api/chirps", apiConfig.handleChirp)

	mux.HandleFunc("POST /api/users", apiConfig.handleUsers)
	mux.HandleFunc("PUT /api/users", apiConfig.handleUpdateUsers)

	mux.HandleFunc("POST /api/login", apiConfig.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiConfig.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handleRevokeToken)

	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlePolkaWebhooks)

	server := &http.Server{Addr: ":" + port, Handler: mux}

	fmt.Printf("⚡  Serving files from %s on port: %s...\n", filepathRoot, port)

	log.Fatal(server.ListenAndServe())
}

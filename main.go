package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	database "github.com/sainikmandal/chirpy-go-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	PolkaKey       string
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is required")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dbQueries, err := database.Prepare(context.Background(), conn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbQueries.Close()

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		PolkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.userCreateHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpCreateHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.chirpGetHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpGetByIDHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebhookHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.userUpdateHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.chirpDeleteHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
    <form action="/admin/reset" method="POST">
      <button type="submit">Reset Counter</button>
    </form>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

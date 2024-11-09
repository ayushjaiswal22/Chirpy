package main

import (
    "fmt"	
    "os"	
    "log"
	"net/http"
    "database/sql"
    "sync/atomic"
    "github.com/ayushjaiswal22/chirpy/internal/database"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
    Db *database.Queries
    Platform string
    SecretKey string
}

var notAllowed map[string]bool

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

var cfg apiConfig

func main() {
	const filepathRoot = "."
	const port = "8080"
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
    fmt.Println(dbURL)
    db, err := sql.Open("postgres", dbURL)
    if err!=nil {
        log.Fatal(err)
    }
    dbQueries := database.New(db)
    cfg = apiConfig{Db:dbQueries, Platform:os.Getenv("PLATFORM"), SecretKey:os.Getenv("SECRET_KEY")}
    cfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", handlerHits)
	mux.HandleFunc("POST /api/users", handlerCreateUser)
	mux.HandleFunc("POST /api/login", handlerLoginUser)
	mux.HandleFunc("POST /admin/reset", handlerResetUsers)
	mux.HandleFunc("POST /api/chirps", handlerChirp)
	mux.HandleFunc("GET /api/chirps", handlerGetAllChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpById)
	mux.HandleFunc("POST /api/refresh", handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", handlerRevokeToken)
	mux.HandleFunc("PUT /api/users", handlerUpdateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


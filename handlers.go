package main

import (
    "strconv"
	"net/http"
)



func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
    hits := strconv.Itoa(int(cfg.fileserverHits.Load()))

    w.Write([]byte(`<html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited ` + hits + ` times!</p>
        </body>
    </html>`))
	w.WriteHeader(http.StatusOK)
}

func handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


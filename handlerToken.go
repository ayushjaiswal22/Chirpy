package main

import (
    "net/http"
    "encoding/json"
    "time"
    "database/sql"
    "github.com/ayushjaiswal22/chirpy/internal/auth"
    "github.com/ayushjaiswal22/chirpy/internal/database"
)


type TokenResponse struct {
    Token string `json:"token"` 
}

func handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
    tokenString, err := auth.GetBearerToken(r.Header)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }
    tokenParams := database.GetUserByTokenParams{Token:tokenString, ExpiresAt:time.Now().UTC()}
    userId, err := cfg.Db.GetUserByToken(r.Context(), tokenParams)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Unable to create Refresh Token\"}"))
        return
    }
    access_token, err := auth.MakeJWT(userId, cfg.SecretKey, time.Hour)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Unable to create Access Token\"}"))
        return
    }
    

    w.WriteHeader(http.StatusOK)
    tokenResp := TokenResponse{
        Token: access_token,
    }
    
    data, err := json.Marshal(tokenResp)
    if err!=nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Write([]byte(data))

}




func handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

    tokenString, err := auth.GetBearerToken(r.Header)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }
    t := sql.NullTime{Time:time.Now().UTC(), Valid:true}
    revokeParams := database.RevokeTokenParams{Token:tokenString, RevokedAt:t}
    err = cfg.Db.RevokeToken(r.Context(), revokeParams)
    if err!=nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

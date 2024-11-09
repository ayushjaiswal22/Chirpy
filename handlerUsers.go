package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "time"
    "strconv"
    "github.com/ayushjaiswal22/chirpy/internal/database"
    "github.com/ayushjaiswal22/chirpy/internal/auth"
)


type CreateUserRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
    ExpiresInSeconds string `json:"expires_in_seconds"`
}

type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"` 
    Email     string `json:"email"`
    AccessToken     string `json:"token"`
    RefreshToken     string `json:"refresh_token"`
}




func handlerCreateUser(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var js CreateUserRequest
    err := decoder.Decode(&js)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }
    fmt.Println(js.Email)
    t := time.Now().UTC()
    passwd, err := auth.HashPassword(js.Password)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return
    }
    args := database.CreateUserParams {
        ID: uuid.New(),
        CreatedAt: t,
        UpdatedAt: t,
        Email: js.Email,
        HashedPassword: passwd,
    }
    u, er := cfg.Db.CreateUser(r.Context(), args)
    if er!=nil {
        fmt.Println(er)
        http.Error(w, "Internal server error 1", http.StatusInternalServerError)
        return
    }
    user := User{
		ID:        u.ID,
        CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
        AccessToken: "",
        RefreshToken: "",
     }
    
    w.WriteHeader(http.StatusCreated)
    data, err := json.Marshal(user)
    if err!=nil {
        http.Error(w, "Internal server error 2", http.StatusInternalServerError)
        return
    }
    w.Write([]byte(data))

    
}

func handlerResetUsers(w http.ResponseWriter, r *http.Request) {
    if cfg.Platform != "dev" {
        http.Error(w, "Forbidden: You do not have permission to access this resource.", http.StatusForbidden)
        return
    }
    w.WriteHeader(http.StatusOK)
    cfg.Db.DeleteUsers(r.Context())
    
}


func handlerLoginUser(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var userReq CreateUserRequest
    err := decoder.Decode(&userReq)
    if err!=nil { 
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return
    }

    userRes, err := cfg.Db.GetHashedPassword(r.Context(), userReq.Email)
    if err!=nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    fmt.Println(userReq.Password)
    fmt.Println(userRes.HashedPassword)
    err = auth.CheckPasswordHash(userReq.Password, userRes.HashedPassword)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Incorrect email or password\"}"))
        return
    }
    var access_token string
    if userReq.ExpiresInSeconds != "" {
        d, _ := strconv.Atoi(userReq.ExpiresInSeconds)
        access_token, err = auth.MakeJWT(userRes.ID, cfg.SecretKey, time.Second * time.Duration(d))
        if err!=nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("{\"error\":\"Unable to create Access Token\"}"))
            return
        }
    } else {
        access_token, err = auth.MakeJWT(userRes.ID, cfg.SecretKey, time.Hour)
        if err!=nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("{\"error\":\"Unable to create Access Token\"}"))
            return
        }

    }

    refresh_token, err := auth.MakeRefreshToken()
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Unable to create Refresh Token while Login\"}"))
        return
    }
    t := time.Now().UTC()
    expire := t.Add(time.Hour * 24 * 60)
    createTokenParams := database.CreateRefreshTokenParams{
        Token: refresh_token,
        CreatedAt: t,
        UserID: userRes.ID,
        ExpiresAt: expire,
    }
    _, err = cfg.Db.CreateRefreshToken(r.Context(), createTokenParams)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Unable to register Refresh Token while Login\"}"))
        return
    }
    w.WriteHeader(http.StatusOK)
    user := User{
        ID:        userRes.ID,
        CreatedAt: userRes.CreatedAt,
        UpdatedAt: userRes.UpdatedAt,
        Email:     userRes.Email,
        AccessToken:     access_token,
        RefreshToken:    refresh_token,
    }


    data, err := json.Marshal(user)
    if err!=nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(data))

}


func handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
    
    decoder := json.NewDecoder(r.Body)
    var userReq CreateUserRequest
    err := decoder.Decode(&userReq)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }
    fmt.Println(userReq.Email)
    hashedPswd, err := auth.HashPassword(userReq.Password)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return
    }

    tokenString, err := auth.GetBearerToken(r.Header)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Could not get Bearer Token\"}"))
        return
    }

    userId, err := auth.ValidateJWT(tokenString, cfg.SecretKey)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"401 Unauthorized\"}"))
        return
    }
    
    userParams := database.UpdateUserParams{
        ID: userId,
        Email: userReq.Email,
        HashedPassword: hashedPswd,
        UpdatedAt: time.Now().UTC(),
    }
    u, err := cfg.Db.UpdateUser(r.Context(), userParams)
    if err!=nil {
        fmt.Println(err)
        http.Error(w, "Internal server error 1", http.StatusInternalServerError)
        return
    }

    user := User{
		ID:        u.ID,
        CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
        AccessToken: "",
        RefreshToken: "",
     }
    
    w.WriteHeader(http.StatusOK)
    data, err := json.Marshal(user)
    if err!=nil {
        http.Error(w, "Internal server error 2", http.StatusInternalServerError)
        return
    }
    w.Write([]byte(data))
}

package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "time"
    "github.com/ayushjaiswal22/chirpy/internal/database"
)


type CreateUserRequest struct {
    Email string `json:"email"`
}

type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"` 
    Email     string `json:"email"`
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
    args := database.CreateUserParams {
        ID: uuid.New(),
        CreatedAt: t,
        UpdatedAt: t,
        Email: js.Email,
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


package main

import (
    "net/http"
    "strings"
    "encoding/json"
    "github.com/google/uuid"
    "time"
    "github.com/ayushjaiswal22/chirpy/internal/database"
    "github.com/ayushjaiswal22/chirpy/internal/auth"
)



func validateChirp(chirp *string) bool {
    if len(*chirp)>140 {
        return false
    }
    notAllowed := map[string]bool {
        "kerfuffle":true,
        "sharbert":true,
        "fornax":true,
    }
    input := strings.Split(*chirp, " ")
    for i, str := range input {
        str = strings.ToLower(str)
        if notAllowed[str] {
            input[i] = "****"
        }
    }
    *chirp = strings.Join(input, " ")
    return true
}


type ChirpRequest struct {
    Body string `json:"body"`
    UserID string `json:"user_id"`
}

type Chirp struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}


func handlerChirp(w http.ResponseWriter, r *http.Request) {    
    w.Header().Add("Content-Type", "application/json; charset=utf-8")
    decoder := json.NewDecoder(r.Body)
    var chirpReq ChirpRequest
    err := decoder.Decode(&chirpReq)
    if err!=nil{
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong while decoding\"}"))
        return 
    }
    tokenString, err := auth.GetBearerToken(r.Header)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"Could not get Bearer Token\"}"))
        return
    }
    uid, err := auth.ValidateJWT(tokenString, cfg.SecretKey)
    if err!=nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("{\"error\":\"401 Unauthorized\"}"))
        return
    }
    if validateChirp(&chirpReq.Body) {
        t := time.Now().UTC()
        chirpRequest := database.CreateChirpParams{
            ID: uuid.New(),
            CreatedAt: t,
            UpdatedAt: t,
            Body:chirpReq.Body,
            UserID:uid,
        }
        chirpResp, err := cfg.Db.CreateChirp(r.Context(), chirpRequest)
        if err!=nil {
            http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        chirp := Chirp {
            ID: chirpResp.ID,
            CreatedAt: chirpResp.CreatedAt,
            UpdatedAt: chirpResp.UpdatedAt,
            Body: chirpResp.Body,
            UserID: chirpResp.UserID,
        }

        data, err := json.Marshal(chirp)
        if err!=nil {
            http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
            return
        }
        w.Write([]byte(data))
    } else {
        http.Error(w, "Not a valid chirp", http.StatusBadRequest)
    }
}

func handlerGetAllChirp(w http.ResponseWriter, r *http.Request) {
    authorID := r.URL.Query().Get("author_id")
    sort := r.URL.Query().Get("sort")
    var chirpsResp []database.Chirp
    var er error
    if authorID=="" {
        if sort=="" || sort=="asc" {
            chirpsResp, er = cfg.Db.GetAllChirps(r.Context())
        } else {
            chirpsResp, er = cfg.Db.GetAllChirpsDesc(r.Context()) 
        }
    } else {
        parsedUserId, err := uuid.Parse(authorID)
        if err!=nil {
            http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
            return
        }
        if sort=="" || sort=="asc"{
            chirpsResp, er = cfg.Db.GetChirpsByUser(r.Context(), parsedUserId)
        } else {
            chirpsResp, er = cfg.Db.GetChirpsByUserDesc(r.Context(), parsedUserId)
        }
    }
    if er!=nil {
        http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
        return
    }

    chirps := make([]Chirp, len(chirpsResp))
    for i, chirp := range chirpsResp {
        tmp := Chirp {
            ID: chirp.ID,
            CreatedAt: chirp.CreatedAt,
            UpdatedAt: chirp.UpdatedAt,
            Body: chirp.Body,
            UserID: chirp.UserID,
        }
        chirps[i] = tmp
    } 
    data, err := json.Marshal(chirps)
    if err!=nil {
        http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
        return
    }
    w.Write([]byte(data))
}


func handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
    chirpId := r.PathValue("chirpID")
    parseChirpId, err := uuid.Parse(chirpId)
    if err!=nil{
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }

    chirpResp, err := cfg.Db.GetChirpById(r.Context(), parseChirpId)
    if err!=nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("{\"error\":\"404 Chirp not found.\"}"))
        return
    } 
    chirp := Chirp {
        ID: chirpResp.ID,
        CreatedAt: chirpResp.CreatedAt,
        UpdatedAt: chirpResp.UpdatedAt,
        Body: chirpResp.Body,
        UserID: chirpResp.UserID,
    } 
    data, err := json.Marshal(chirp)
    if err!=nil {
        http.Error(w, "Internal Server error while chirping 1", http.StatusInternalServerError)
        return
    }
    w.Write([]byte(data))
}



func handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
   chirpId := r.PathValue("chirpID")
   parseChirpId, err := uuid.Parse(chirpId)

    if err!=nil{
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
    uid, err := auth.ValidateJWT(tokenString, cfg.SecretKey)
    if err!=nil {
        w.WriteHeader(http.StatusForbidden)
        w.Write([]byte("{\"error\":\"403 Forbidden\"}"))
        return
    }
    chirpResp, err := cfg.Db.GetChirpById(r.Context(), parseChirpId)
    if uid != chirpResp.UserID {
        w.WriteHeader(http.StatusForbidden)
        w.Write([]byte("{\"error\":\"403 Forbidden\"}"))
        return
    }
    
    err = cfg.Db.DeleteChirpById(r.Context(), parseChirpId)
    if err!=nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("{\"error\":\"404 Chirp not found.\"}"))
        return
    } 
    w.WriteHeader(http.StatusNoContent)
}


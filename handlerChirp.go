package main

import (
    "net/http"
    "strings"
    "encoding/json"
    "github.com/google/uuid"
    "time"
    "github.com/ayushjaiswal22/chirpy/internal/database"
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
        w.Write([]byte("{\"error\":\"Something went wrong\"}"))
        return 
    }
    if validateChirp(&chirpReq.Body) {
        parsedUserID, err := uuid.Parse(chirpReq.UserID)
        if err!=nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("{\"error\":\"Something went wrong\"}"))
            return 
        }

        t := time.Now().UTC()
        chirpRequest := database.CreateChirpParams{
            ID: uuid.New(),
            CreatedAt: t,
            UpdatedAt: t,
            Body:chirpReq.Body,
            UserID:parsedUserID,
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
    chirpsResp, err := cfg.Db.GetAllChirps(r.Context())
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

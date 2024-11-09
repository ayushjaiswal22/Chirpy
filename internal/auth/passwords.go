package auth

import (
    "github.com/google/uuid"
    "net/http"
    "fmt"
    "time"
    "strings"
    "errors"
    "crypto/rand"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "encoding/hex"
)


func MakeRefreshToken() (string, error) {
    data := make([]byte, 32)
    _, err := rand.Read(data)
    if err!=nil {
        return "", nil
    }
    refreshToken := hex.EncodeToString(data)
    return refreshToken, nil
}

type MyCustomClaims struct {
    
    jwt.RegisteredClaims
}
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString([]byte(tokenSecret))
    return ss, err

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    } else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
        fmt.Println(claims.Subject, claims.Issuer)
        parsedUUID, err := uuid.Parse(claims.Subject)
        if err!=nil {
            return uuid.Nil, nil
        }
        return parsedUUID, nil
    } else {
        return uuid.Nil, err
    }
}

func GetBearerToken(headers http.Header) (string, error) {
    tokenString := headers.Get("Authorization")
    if strings.Contains(tokenString, "Bearer ") {
        tokenString = tokenString[len("Bearer "):]
        return tokenString, nil
    } else {
        return "", errors.New("Malformed access token")
    }
}

func HashPassword(password string) (string, error) {
    data, err := bcrypt.GenerateFromPassword([]byte(password), 5)
    return string(data), err
}

func CheckPasswordHash(password, hash string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err
}

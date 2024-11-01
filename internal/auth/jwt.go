package auth

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("votre_clé_secrète_super_sécurisée")

type Claims struct {
    UserID  uint   `json:"UserId"`
    LoginID string `json:"loginId"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, loginID string) (string, error) {
    claims := &Claims{
        UserID:  userID,
        LoginID: loginID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("token invalide")
    }

    return claims, nil
}

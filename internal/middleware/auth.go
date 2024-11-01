package middleware

import (
    "context"
    "net/http"
    "strings"
    "fmt"
    "Projet_soul_back/internal/auth"
)

type contextKey string
const UserIDKey contextKey = "UserID"  // Majuscule ici aussi

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            next.ServeHTTP(w, r)
            return
        }

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header requis", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            http.Error(w, "Format d'autorisation invalide", http.StatusUnauthorized)
            return
        }

        claims, err := auth.ValidateToken(tokenParts[1])
        if err != nil {
            http.Error(w, "Token invalide: "+err.Error(), http.StatusUnauthorized)
            return
        }

        if claims.UserID == 0 {
            http.Error(w, "ID utilisateur invalide dans le token", http.StatusUnauthorized)
            return
        }

        // Ajout de debug logging
        fmt.Printf("User ID from token: %d\n", claims.UserID)

        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

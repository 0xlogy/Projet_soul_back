package middleware

import (
    "fmt"
    "net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log les requêtes entrantes pour debug
        fmt.Printf("Requête entrante: %s %s\n", r.Method, r.URL.Path)
        fmt.Printf("Headers origine: %v\n", r.Header.Get("Origin"))

        // Configuration des headers CORS
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "*")
        w.Header().Set("Access-Control-Max-Age", "86400")
        w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Gestion de la requête OPTIONS (preflight)
        if r.Method == "OPTIONS" {
            fmt.Println("Requête OPTIONS reçue")
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

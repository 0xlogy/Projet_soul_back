package api

import (
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/api/handlers"
    "Projet_soul_back/internal/middleware"
)

// CorsMiddleware gère les headers CORS
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Ajout des headers CORS plus permissifs
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
        w.Header().Set("Access-Control-Max-Age", "3600")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        // Gestion de la requête OPTIONS préliminaire
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    
    r.Use(corsMiddleware) 
    // Route publique pour l'authentification
    r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

    // Routes protégées
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)

    // Routes utilisateurs
    api.HandleFunc("/users", handlers.GetUsers).Methods("GET", "OPTIONS")
    api.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET", "OPTIONS")
    api.HandleFunc("/users/{id}/souls", handlers.UpdateUserSouls).Methods("PATCH", "OPTIONS")

    // Routes jeux
    api.HandleFunc("/games", handlers.GetGames).Methods("GET")
    api.HandleFunc("/games/{id}", handlers.GetGame).Methods("GET")

    // Routes scores
    api.HandleFunc("/scores", handlers.GetScores).Methods("GET")
    api.HandleFunc("/scores", handlers.CreateScore).Methods("POST", "OPTIONS")
    api.HandleFunc("/users/{userId}/scores", handlers.GetUserScores).Methods("GET")

    // Routes statistiques spécifiques aux jeux
    api.HandleFunc("/games/chisoulmi/stats", handlers.GetChiSoulMiStats).Methods("GET", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/versus/{player1Id}/{player2Id}", handlers.GetChiSoulMiVersusStats).Methods("GET", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/leaderboard", handlers.GetChiSoulMiLeaderboard).Methods("GET", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/duels", handlers.CreateDuel).Methods("POST", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/duels/join", handlers.JoinDuel).Methods("POST", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/duels/{duelId}", handlers.GetDuel).Methods("GET", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/duels/{duelId}/choice", handlers.MakeChoice).Methods("POST", "OPTIONS")
    api.HandleFunc("/games/chisoulmi/duels/{duelId}/timeout", handlers.TimeoutDuel).Methods("POST", "OPTIONS")
    api.HandleFunc("/games/pacman/stats", handlers.GetPacmanStats).Methods("GET")

    // Route admin pour la mise à jour des âmes (à protéger avec un middleware admin)
    api.HandleFunc("/admin/users/{id}/souls", handlers.AdminUpdateSouls).Methods("PATCH")


    // Appliquer le middleware CORS à toutes les routes
    r.Use(corsMiddleware)

    return r
}

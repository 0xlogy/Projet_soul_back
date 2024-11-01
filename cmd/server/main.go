package main

import (
    "fmt"
    "log"
    "github.com/gorilla/mux"
    "net/http"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/middleware"
    "Projet_soul_back/internal/api"
)

func main() {
    // Initialisation de la base de données
    db, err := database.InitDB()
    if err != nil {
        log.Fatalf("Erreur d'initialisation de la base de données: %v", err)
    }

    // Initialisation des données par défaut
    if err := database.InitializeDB(db); err != nil {
        log.Fatalf("Erreur d'initialisation des données: %v", err)
    }

    // Configuration du routeur
    router := api.SetupRouter()

    // Ajout du middleware CORS
    router.Use(middleware.CorsMiddleware)

    // Démarrage du serveur
    port := ":8080"
    fmt.Printf("Serveur démarré sur le port %s\n", port)
    fmt.Println("Routes configurées:")
    router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        path, _ := route.GetPathTemplate()
        methods, _ := route.GetMethods()
        fmt.Printf("Route: %s [%v]\n", path, methods)
        return nil
    })

    if err := http.ListenAndServe(port, router); err != nil {
        log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
    }
}

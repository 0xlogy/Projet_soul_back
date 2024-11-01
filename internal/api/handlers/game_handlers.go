package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

// GetGames retourne la liste de tous les jeux
func GetGames(w http.ResponseWriter, r *http.Request) {
    var games []models.Game
    if err := database.DB.Find(&games).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(games)
}

// GetGame retourne un jeu spécifique
func GetGame(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var game models.Game
    if err := database.DB.First(&game, id).Error; err != nil {
        http.Error(w, "Jeu non trouvé", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(game)
}

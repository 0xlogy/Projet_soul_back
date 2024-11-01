package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

// GetScores retourne tous les scores
func GetScores(w http.ResponseWriter, r *http.Request) {
    var scores []models.Score
    if err := database.DB.Find(&scores).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(scores)
}

// CreateScore crée un nouveau score
func CreateScore(w http.ResponseWriter, r *http.Request) {
    var score models.Score
    if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&score).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(score)
}

// GetUserScores retourne les scores d'un utilisateur spécifique
func GetUserScores(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["userId"]

    var scores []models.Score
    if err := database.DB.Where("user_id = ?", userId).Find(&scores).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(scores)
}

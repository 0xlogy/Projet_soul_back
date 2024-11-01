// internal/api/handlers/user_handlers.go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

// GetUsers retourne la liste de tous les utilisateurs
func GetUsers(w http.ResponseWriter, r *http.Request) {
    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// GetUser retourne un utilisateur spécifique
func GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// UpdateUserSouls met à jour le nombre d'âmes d'un utilisateur
func UpdateUserSouls(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var updateData struct {
        Souls int `json:"souls"`
    }

    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
        return
    }

    user.Souls = updateData.Souls
    if err := database.DB.Save(&user).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

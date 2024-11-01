package handlers

import (
    "encoding/json"
    "net/http"
    "Projet_soul_back/internal/auth"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

type LoginRequest struct {
    LoginID  string `json:"loginId"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User  models.User `json:"user"`
    Token string      `json:"token"`
}

// Login gère l'authentification des utilisateurs
func Login(w http.ResponseWriter, r *http.Request) {
    var loginReq LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := database.DB.Where("login_id = ?", loginReq.LoginID).First(&user).Error; err != nil {
        http.Error(w, "Identifiant ou mot de passe incorrect", http.StatusUnauthorized)
        return
    }

    if !auth.CheckPasswordHash(loginReq.Password, user.Password) {
        http.Error(w, "Identifiant ou mot de passe incorrect", http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateToken(user.ID, user.LoginID)
    if err != nil {
        http.Error(w, "Erreur lors de la création du token", http.StatusInternalServerError)
        return
    }

    response := LoginResponse{
        User:  user,
        Token: token,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

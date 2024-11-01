package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

// ChiSoulMiStats représente les statistiques du Chi-Soul-Mi pour un joueur
type ChiSoulMiStats struct {
    PlayerID    uint   `json:"playerId"`
    PlayerName  string `json:"playerName"`
    TotalWins   int    `json:"totalWins"`
    Opponents   []OpponentStats `json:"opponents"`
}

type OpponentStats struct {
    OpponentID   uint   `json:"opponentId"`
    OpponentName string `json:"opponentName"`
    WinsAgainst  int    `json:"winsAgainst"`
}

// PacmanStats représente les statistiques de Pacman
type PacmanStats struct {
    PlayerID    uint   `json:"playerId"`
    PlayerName  string `json:"playerName"`
    BestScore   int    `json:"bestScore"`
}

// GetChiSoulMiStats retourne les statistiques de Chi-Soul-Mi pour un joueur
func GetChiSoulMiStats(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userId").(uint)

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
        return
    }

    var stats ChiSoulMiStats
    stats.PlayerID = user.ID
    stats.PlayerName = user.Nickname

    // Récupère le jeu Chi-Soul-Mi
    var chiSoulMi models.Game
    if err := database.DB.Where("name = ?", "Chi-Soul-Mi").First(&chiSoulMi).Error; err != nil {
        http.Error(w, "Jeu non trouvé", http.StatusNotFound)
        return
    }

    // Compte le nombre total de victoires
    var totalWins int64
    database.DB.Model(&models.Score{}).
        Where("user_id = ? AND game_id = ? AND result = 'win'", userID, chiSoulMi.ID).
        Count(&totalWins)
    stats.TotalWins = int(totalWins)

    // Récupère les statistiques par adversaire
    rows, err := database.DB.Raw(`
        SELECT 
            opponent_id,
            u.nickname as opponent_name,
            COUNT(*) as wins_against
        FROM scores s
        JOIN users u ON u.id = s.opponent_id
        WHERE s.user_id = ? 
        AND s.game_id = ?
        AND s.result = 'win'
        GROUP BY opponent_id, u.nickname`, userID, chiSoulMi.ID).Rows()

    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var opStats OpponentStats
            rows.Scan(&opStats.OpponentID, &opStats.OpponentName, &opStats.WinsAgainst)
            stats.Opponents = append(stats.Opponents, opStats)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// GetPacmanStats retourne le meilleur score à Pacman pour un joueur
func GetPacmanStats(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userId").(uint)

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
        return
    }

    var stats PacmanStats
    stats.PlayerID = user.ID
    stats.PlayerName = user.Nickname

    // Récupère le jeu Pacman
    var pacman models.Game
    if err := database.DB.Where("name = ?", "Pacman").First(&pacman).Error; err != nil {
        http.Error(w, "Jeu non trouvé", http.StatusNotFound)
        return
    }

    // Récupère le meilleur score
    var bestScore models.Score
    if err := database.DB.Where("user_id = ? AND game_id = ?", userID, pacman.ID).
        Order("points DESC").
        First(&bestScore).Error; err == nil {
        stats.BestScore = bestScore.Points
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// SoulsUpdateRequest est la structure pour la mise à jour des âmes
type SoulsUpdateRequest struct {
    Souls int `json:"souls"`
}

// AdminUpdateSouls est un endpoint réservé aux administrateurs
func AdminUpdateSouls(w http.ResponseWriter, r *http.Request) {
    // Vérifier si l'utilisateur est admin (à implémenter selon tes besoins)
    isAdmin := checkIfAdmin(r)
    if !isAdmin {
        http.Error(w, "Accès non autorisé", http.StatusForbidden)
        return
    }

    var updateReq SoulsUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    // Utilisation de vars["id"] au lieu de redéclarer userID
    vars := mux.Vars(r)
    targetUserID := vars["id"]

    var user models.User
    if err := database.DB.First(&user, targetUserID).Error; err != nil {
        http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
        return
    }

    user.Souls = updateReq.Souls
    if err := database.DB.Save(&user).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Fonction helper pour vérifier si l'utilisateur est admin
func checkIfAdmin(r *http.Request) bool {
    // À implémenter selon ta logique d'administration
    // Par exemple, vérifier un rôle dans le token JWT
    return true // Pour l'instant, retourne toujours true
}

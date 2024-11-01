package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
)

// VersusStats représente les statistiques entre deux joueurs
type VersusStats struct {
    Player1 struct {
        ID       uint   `json:"id"`
        Nickname string `json:"nickname"`
        Wins     int    `json:"wins"`
    } `json:"player1"`
    Player2 struct {
        ID       uint   `json:"id"`
        Nickname string `json:"nickname"`
        Wins     int    `json:"wins"`
    } `json:"player2"`
    TotalGames int `json:"totalGames"`
}

// GetChiSoulMiVersusStats retourne les statistiques entre deux joueurs
func GetChiSoulMiVersusStats(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player1ID := vars["player1Id"]
    player2ID := vars["player2Id"]

    var stats VersusStats

    // Récupère le jeu Chi-Soul-Mi
    var chiSoulMi models.Game
    if err := database.DB.Where("name = ?", "Chi-Soul-Mi").First(&chiSoulMi).Error; err != nil {
        http.Error(w, "Jeu non trouvé", http.StatusNotFound)
        return
    }

    // Récupère les informations des joueurs
    var player1, player2 models.User
    if err := database.DB.First(&player1, player1ID).Error; err != nil {
        http.Error(w, "Joueur 1 non trouvé", http.StatusNotFound)
        return
    }
    if err := database.DB.First(&player2, player2ID).Error; err != nil {
        http.Error(w, "Joueur 2 non trouvé", http.StatusNotFound)
        return
    }

    stats.Player1.ID = player1.ID
    stats.Player1.Nickname = player1.Nickname
    stats.Player2.ID = player2.ID
    stats.Player2.Nickname = player2.Nickname

    // Compte les victoires de player1 contre player2
    var wins1 int64
    database.DB.Model(&models.Score{}).
        Where("game_id = ? AND user_id = ? AND opponent_id = ? AND result = 'win'", 
              chiSoulMi.ID, player1.ID, player2.ID).
        Count(&wins1)
    stats.Player1.Wins = int(wins1)

    // Compte les victoires de player2 contre player1
    var wins2 int64
    database.DB.Model(&models.Score{}).
        Where("game_id = ? AND user_id = ? AND opponent_id = ? AND result = 'win'", 
              chiSoulMi.ID, player2.ID, player1.ID).
        Count(&wins2)
    stats.Player2.Wins = int(wins2)

    // Compte le nombre total de parties entre les deux joueurs
    var total int64
    database.DB.Model(&models.Score{}).
        Where("game_id = ? AND ((user_id = ? AND opponent_id = ?) OR (user_id = ? AND opponent_id = ?))",
              chiSoulMi.ID, player1.ID, player2.ID, player2.ID, player1.ID).
        Count(&total)
    stats.TotalGames = int(total)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// GetChiSoulMiLeaderboard retourne le classement global du Chi-Soul-Mi
func GetChiSoulMiLeaderboard(w http.ResponseWriter, r *http.Request) {
    type PlayerStats struct {
        ID       uint    `json:"id"`
        Nickname string  `json:"nickname"`
        Wins     int     `json:"wins"`
        Losses   int     `json:"losses"`
        WinRate  float64 `json:"winRate"`
    }

    var chiSoulMi models.Game
    if err := database.DB.Where("name = ?", "Chi-Soul-Mi").First(&chiSoulMi).Error; err != nil {
        http.Error(w, "Jeu non trouvé", http.StatusNotFound)
        return
    }

    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var leaderboard []PlayerStats
    for _, user := range users {
        var stats PlayerStats
        stats.ID = user.ID
        stats.Nickname = user.Nickname

        // Compte les victoires
        var wins int64
        database.DB.Model(&models.Score{}).
            Where("game_id = ? AND user_id = ? AND result = 'win'", chiSoulMi.ID, user.ID).
            Count(&wins)
        stats.Wins = int(wins)

        // Compte les défaites
        var losses int64
        database.DB.Model(&models.Score{}).
            Where("game_id = ? AND user_id = ? AND result = 'lose'", chiSoulMi.ID, user.ID).
            Count(&losses)
        stats.Losses = int(losses)

        totalGames := stats.Wins + stats.Losses
        if totalGames > 0 {
            stats.WinRate = float64(stats.Wins) / float64(totalGames) * 100
        }

        leaderboard = append(leaderboard, stats)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(leaderboard)
}

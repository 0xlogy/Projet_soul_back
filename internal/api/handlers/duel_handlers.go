package handlers

import (
    "fmt"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "Projet_soul_back/internal/database"
    "Projet_soul_back/internal/models"
    "Projet_soul_back/internal/middleware"
    "time"
    "crypto/rand"
    "encoding/base64"
)

func generateInviteCode() (string, error) {
    // Génère 6 bytes aléatoires
    bytes := make([]byte, 6)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    
    // Convertit en base64 et nettoie le résultat
    code := base64.URLEncoding.EncodeToString(bytes)
    // Prend seulement les 8 premiers caractères pour avoir un code court
    code = code[:8]
    
    return code, nil
}  

// GetDuel récupère un duel par son ID
func GetDuel(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    duelID := vars["duelId"]

    var duel models.Duel
    if err := database.DB.First(&duel, duelID).Error; err != nil {
        http.Error(w, "Duel non trouvé", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(duel)
}

// CreateDuel crée un nouveau duel

func CreateDuel(w http.ResponseWriter, r *http.Request) {
    // Activation des logs détaillés
    fmt.Printf("Création d'un nouveau duel\n")

    // Récupération de l'ID utilisateur depuis le contexte
    userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
    if !ok {
        fmt.Printf("Erreur: ID utilisateur non trouvé dans le contexte\n")
        http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
        return
    }

    // Configuration des headers de réponse en premier
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    // Vérification des duels existants
    var existingDuel models.Duel
    if err := database.DB.Where("creator_id = ? AND status IN ('pending', 'active')", userID).First(&existingDuel).Error; err == nil {
        fmt.Printf("Duel existant trouvé pour l'utilisateur %d\n", userID)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Vous avez déjà un duel en cours",
        })
        return
    }

    // Récupération du jeu Chi-Soul-Mi
    var game models.Game
    if err := database.DB.Where("name = ?", "Chi-Soul-Mi").First(&game).Error; err != nil {
        fmt.Printf("Erreur: Jeu Chi-Soul-Mi non trouvé\n")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Jeu non trouvé",
        })
        return
    }

    // Génération du code d'invitation
    inviteCode, err := generateInviteCode()
    if err != nil {
        fmt.Printf("Erreur de génération du code: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Erreur lors de la génération du code",
        })
        return
    }

    // Création du duel
    duel := models.Duel{
        InviteCode: inviteCode,
        CreatorID:  userID,
        GameID:     game.ID,
        Status:     "pending",
        ExpiresAt:  time.Now().Add(15 * time.Minute),
    }

    if err := database.DB.Create(&duel).Error; err != nil {
        fmt.Printf("Erreur de création du duel: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
    }

    fmt.Printf("Duel créé avec succès: %+v\n", duel)

    // Assurons-nous d'avoir un json valide en retour
    response := map[string]interface{}{
        "id":         duel.ID,
        "inviteCode": duel.InviteCode,
        "creatorId":  duel.CreatorID,
        "status":     duel.Status,
        "expiresAt":  duel.ExpiresAt,
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        fmt.Printf("Erreur d'encodage JSON: %v\n", err)
        http.Error(w, "Erreur serveur", http.StatusInternalServerError)
        return
    }
}

// JoinDuel permet à un joueur de rejoindre un duel existant
func JoinDuel(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
    if !ok {
        http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
        return
    }

    var joinRequest struct {
        InviteCode string `json:"inviteCode"`
    }

    if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    var duel models.Duel
    if err := database.DB.Where("invite_code = ?", joinRequest.InviteCode).First(&duel).Error; err != nil {
        http.Error(w, "Duel non trouvé", http.StatusNotFound)
        return
    }

    if duel.Status != "pending" {
        http.Error(w, "Ce duel n'est plus disponible", http.StatusBadRequest)
        return
    }

    if duel.CreatorID == userID {
        http.Error(w, "Vous ne pouvez pas rejoindre votre propre duel", http.StatusBadRequest)
        return
    }

    duel.OpponentID = &userID
    duel.Status = "active"  // Important : mise à jour du statut

    if err := database.DB.Save(&duel).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(duel)
}

// MakeChoice enregistre le choix d'un joueur dans un duel
func MakeChoice(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
    if !ok {
        http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    duelID := vars["duelId"]

    var choiceRequest struct {
        Choice string `json:"choice"`
    }

    if err := json.NewDecoder(r.Body).Decode(&choiceRequest); err != nil {
        http.Error(w, "Format de requête invalide", http.StatusBadRequest)
        return
    }

    var duel models.Duel
    if err := database.DB.First(&duel, duelID).Error; err != nil {
        http.Error(w, "Duel non trouvé", http.StatusNotFound)
        return
    }

    // Vérifie que l'utilisateur fait partie du duel
    if duel.CreatorID != userID && *duel.OpponentID != userID {
        http.Error(w, "Vous ne faites pas partie de ce duel", http.StatusForbidden)
        return
    }

    // Enregistre le choix
    if duel.CreatorID == userID {
        duel.CreatorChoice = &choiceRequest.Choice
    } else {
        duel.OpponentChoice = &choiceRequest.Choice
    }

    // Si les deux joueurs ont fait leur choix, détermine le gagnant
    if duel.CreatorChoice != nil && duel.OpponentChoice != nil {
        winner := determineWinner(*duel.CreatorChoice, *duel.OpponentChoice)
        if winner == "creator" {
            duel.WinnerID = &duel.CreatorID
        } else if winner == "opponent" {
            duel.WinnerID = duel.OpponentID
        }
        duel.Status = "completed"
    }

    if err := database.DB.Save(&duel).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(duel)
}



func determineWinner(choice1, choice2 string) string {
    choices := map[string][]string{
        "scissors": {"paper", "lizard"},
        "paper":    {"rock", "spock"},
        "rock":     {"scissors", "lizard"},
        "lizard":   {"paper", "spock"},
        "spock":    {"scissors", "rock"},
    }

    if choice1 == choice2 {
        return "draw"
    }

    for _, beatenChoice := range choices[choice1] {
        if beatenChoice == choice2 {
            return "creator"
        }
    }

    return "opponent"
}


func TimeoutDuel(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    duelID := vars["duelId"]

    var duel models.Duel
    if err := database.DB.First(&duel, duelID).Error; err != nil {
        http.Error(w, "Duel non trouvé", http.StatusNotFound)
        return
    }

    // Déterminer le gagnant basé sur les scores
    var creatorScore int64
    var opponentScore int64

    // Compter les scores
    err := database.DB.Model(&models.Score{}).
        Where("user_id = ? AND duel_id = ?", duel.CreatorID, duel.ID).
        Count(&creatorScore)
    if err != nil {
        http.Error(w, "Erreur lors du comptage des scores", http.StatusInternalServerError)
        return
    }

    err = database.DB.Model(&models.Score{}).
        Where("user_id = ? AND duel_id = ?", duel.OpponentID, duel.ID).
        Count(&opponentScore)
    if err != nil {
        http.Error(w, "Erreur lors du comptage des scores", http.StatusInternalServerError)
        return
    }

    // Déterminer le gagnant
    var winnerId uint
    if creatorScore > opponentScore {
        winnerId = duel.CreatorID
    } else if opponentScore > creatorScore {
        winnerId = *duel.OpponentID
    } else {
        // En cas d'égalité, on peut soit ne pas mettre de gagnant,
        // soit implémenter une logique spécifique pour les égalités
        winnerId = 0
    }

    // Mettre à jour le duel
    duel.Status = "completed"
    if winnerId != 0 {
        duel.WinnerID = &winnerId
    }

    if err := database.DB.Save(&duel).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Envoyer la réponse
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "completed",
        "winnerId": winnerId,
        "creatorScore": creatorScore,
        "opponentScore": opponentScore,
    })
}

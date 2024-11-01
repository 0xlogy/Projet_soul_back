package database

import (
    "fmt"
    "log"
    "os"
    "Projet_soul_back/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initialise la connexion à la base de données
func InitDB() (*gorm.DB, error) {
    // Essaie d'abord d'utiliser DATABASE_URL
    dbURL := os.Getenv("DATABASE_URL")
    
    // Si DATABASE_URL n'existe pas, construit l'URL manuellement
    if dbURL == "" {
        dbURL = fmt.Sprintf(
            "host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
            os.Getenv("PGHOST"),
            os.Getenv("PGUSER"),
            os.Getenv("PGPASSWORD"),
            "soulsitedb",
            os.Getenv("PGPORT"),
        )
    }

    // Log pour debug (à retirer en production)
    fmt.Printf("Tentative de connexion avec l'URL: %s\n", dbURL)

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("erreur de connexion à la base de données: %v", err)
    }

    DB = db
    return db, nil
}

// InitializeDB configure la base de données avec les données par défaut
func InitializeDB(db *gorm.DB) error {
    // Auto-migration des tables
    if err := db.AutoMigrate(&models.User{}, &models.Game{}, &models.Score{}); err != nil {
        return fmt.Errorf("erreur lors de la migration: %v", err)
    }

    // Création des utilisateurs par défaut
    for _, user := range models.DefaultPlayers {
        var existingUser models.User
        if err := db.Where("nickname = ?", user.Nickname).First(&existingUser).Error; err != nil {
            if err := db.Create(&user).Error; err != nil {
                log.Printf("Erreur lors de la création de l'utilisateur %s: %v", user.Nickname, err)
            } else {
                log.Printf("Utilisateur créé: %s", user.Nickname)
            }
        }
    }

    // Création des jeux par défaut
    for _, game := range models.DefaultGames {
        var existingGame models.Game
        if err := db.Where("name = ?", game.Name).First(&existingGame).Error; err != nil {
            if err := db.Create(&game).Error; err != nil {
                log.Printf("Erreur lors de la création du jeu %s: %v", game.Name, err)
            } else {
                log.Printf("Jeu créé: %s", game.Name)
            }
        }
    }

    return nil
}

// GetDB retourne l'instance de la base de données
func GetDB() *gorm.DB {
    return DB
}

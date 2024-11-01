package models

import (
    "gorm.io/gorm"
)

type Game struct {
    gorm.Model
    Name        string    `json:"name" gorm:"unique;not null"`
    Type        string    `json:"type" gorm:"type:text;check:type IN ('BONUS', 'MALUS');not null"`
    Description string    `json:"description"`
    IsActive    bool      `json:"isActive" gorm:"default:true"`
}

// Configuration initiale des mini-jeux
var DefaultGames = []Game{
    {
        Name: "Quête des âmes",
        Type: "BONUS",
        Description: "Mini-jeu bonus où vous devez choisir intelligemment vos quêtes",
        IsActive: true,
    },
    {
        Name: "Chi-Soul-Mi",
        Type: "MALUS",
        Description: "Version améliorée du Pierre-Papier-Ciseaux avec Lézard et Spock",
        IsActive: true,
    },
    {
        Name: "Les Reliques Perdues",
        Type: "BONUS",
        Description: "Mini-jeu bonus où le temps est votre allié",
        IsActive: true,
    },
    {
        Name: "Écho de l'Âme",
        Type: "BONUS",
        Description: "Mini-jeu bonus de mémoire et de réflexes",
        IsActive: false,
    },
    {
        Name: "Labyrinthe de l'Esprit",
        Type: "MALUS",
        Description: "Mini-jeu malus de navigation dans un labyrinthe",
        IsActive: false,
    },
}

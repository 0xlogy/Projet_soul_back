package models

import (
    "gorm.io/gorm"
)

type Score struct {
    gorm.Model
    UserID      uint      `json:"userId" gorm:"not null"`
    GameID      uint      `json:"gameId" gorm:"not null"`
    Points      int       `json:"points" gorm:"not null"`
    OpponentID  *uint     `json:"opponentId,omitempty"` // Pour Chi-Soul-Mi
    Result      string    `json:"result,omitempty"`// 'win', 'lose', ou 'draw' pour Chi-Soul-Mi
}

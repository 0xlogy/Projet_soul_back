package models

import (
    "gorm.io/gorm"
    "time"
)

type Duel struct {
    gorm.Model
    InviteCode      string    `json:"inviteCode" gorm:"uniqueIndex;not null"`
    CreatorID       uint      `json:"creatorId" gorm:"not null"`
    OpponentID      *uint     `json:"opponentId"`
    GameID          uint      `json:"gameId" gorm:"not null"`
    Status          string    `json:"status" gorm:"type:text;check:status IN ('pending', 'active', 'completed', 'expired');default:'pending'"`
    ExpiresAt       time.Time `json:"expiresAt"`
    CreatorChoice   *string   `json:"creatorChoice,omitempty"`
    OpponentChoice  *string   `json:"opponentChoice,omitempty"`
    WinnerID        *uint     `json:"winnerId,omitempty"`
}

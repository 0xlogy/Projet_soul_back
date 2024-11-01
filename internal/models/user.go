package models

import (
    "Projet_soul_back/internal/auth"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Nickname    string    `json:"nickname" gorm:"unique;not null"`
    LoginID     string    `json:"loginId" gorm:"unique;not null"`
    Password    string    `json:"-" gorm:"not null"`
    Souls       int       `json:"souls" gorm:"default:0"`
    Scores      []Score   `json:"scores,omitempty"`
}

// BeforeCreate est un hook GORM qui hash le mot de passe avant la création
func (u *User) BeforeCreate(tx *gorm.DB) error {
    hashedPassword, err := auth.HashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hashedPassword
    return nil
}

// Configuration initiale des joueurs
var DefaultPlayers = []User{
    {
        Nickname: "Zaza",
        LoginID: "BaronneDuNesquik",
        Password: "GrosseCochonneZaza6996@",
        Souls: 0,
    },
    {
        Nickname: "Charly",
        LoginID: "SaladeDesigner",
        Password: "LeLitConjugalAubergine2552@",
        Souls: 0,
    },
    {
        Nickname: "Gray",
        LoginID: "FiertePlusUltra",
        Password: "MirajaneForeverAndEver5775@",
        Souls: 0,
    },
    {
        Nickname: "Akuma",
        LoginID: "HuitreDeLOmbre",
        Password: "OksanaForeverAndEver1010@",
        Souls: 0,
    },
    {
        Nickname: "Quasibrother",
        LoginID: "BonPartiOuiOui",
        Password: "CecekissForeverAndEverBestPipe1771@",
        Souls: 0,
    },
    {
        Nickname: "Shishi",
        LoginID: "SenseiFantome",
        Password: "BestHealerOfTheEdenEternal9999@",
        Souls: 0,
    },
    {
        Nickname: "Evil Luxus",
        LoginID: "MasterViking",
        Password: "LaNonConnexionDuCul0404@",
        Souls: 0,
    },
    {
        Nickname: "Atsina",
        LoginID: "TelephoneRose",
        Password: "Lalalou95@",
        Souls: 0,
    },
    {
        Nickname: "Billy",
        LoginID: "HighNoon",
        Password: "DesignerQuiDechireDuCul8228@",
        Souls: 0,
    },
    {
        Nickname: "Logy",
        LoginID: "LegendeFoudroye",
        Password: "TheMeeeeelHoooooleeeee0000@",
        Souls: 0,
    },
}

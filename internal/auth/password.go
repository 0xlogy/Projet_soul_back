package auth

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword crée un hash bcrypt d'un mot de passe
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash vérifie si un mot de passe correspond à un hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

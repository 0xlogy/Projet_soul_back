# Build stage
FROM golang:1.22-alpine AS builder

# Installation des dépendances essentielles
RUN apk --no-cache add gcc musl-dev

# Définition du répertoire de travail
WORKDIR /app

# Copie des fichiers de dépendances
COPY go.mod go.sum ./

# Téléchargement des dépendances
RUN go mod download

# Copie du code source
COPY . .

# Construction de l'application avec des optimisations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o main ./cmd/server

# Image finale légère
FROM alpine:latest

# Installation des certificats CA pour HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copie du binaire depuis l'étape de build
COPY --from=builder /app/main .

# Exposition du port (utilisera PORT from env en production)
EXPOSE 8080

# Commande de démarrage
CMD ["./main"]

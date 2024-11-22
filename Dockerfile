FROM golang:1.23 AS build

WORKDIR /app
COPY . .

# Installer ffmpeg et les dépendances Go
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    go mod tidy

# Compiler l'application Go en statique
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Étape de création de l'image finale avec distroless
FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg && \
    rm -rf /var/lib/apt/lists/*
# Copier l'exécutable Go et ffmpeg de l'étape de build
COPY --from=build /app/app /

# Définir le point d'entrée de l'application
CMD ["/app"]

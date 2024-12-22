# Étape 1 : Utiliser une image Go officielle
FROM golang:1.20 as builder

# Installer les dépendances
RUN apt-get update && apt-get install -y \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Télécharger les bibliothèques TensorFlow officielles
RUN wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.11.0.tar.gz && \
    tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-2.11.0.tar.gz && \
    ldconfig && \
    rm libtensorflow-cpu-linux-x86_64-2.11.0.tar.gz

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers nécessaires
COPY . .

# Télécharger les dépendances et construire le projet
RUN go mod tidy && go build -o main .

# Étape 2 : Exécuter l'application dans une image propre
FROM debian:bullseye-slim

# Installer les bibliothèques TensorFlow nécessaires
RUN apt-get update && apt-get install -y \
    libtensorflow2 \
    libtensorflow2-dev \
    && rm -rf /var/lib/apt/lists/*

# Copier le binaire depuis le conteneur de construction
COPY --from=builder /app/main /usr/local/bin/main

# Copier les fichiers nécessaires
COPY cascade /app/cascade
COPY cache /app/cache
COPY testdata /app/testdata

# Définir le point d'entrée
ENTRYPOINT ["main"]

package main

import (
    "encoding/json"
    "flag"
    "log"
    "os"

    "facedetection/internal/face"
)

func main() {
    cascadePath := flag.String("cascade", "cascade.dat", "Chemin vers le fichier cascade Pigo")
    modelPath := flag.String("model", "facenet", "Chemin vers le modèle FaceNet")
    imagePath := flag.String("image", "", "Chemin vers l'image à analyser")
    outputPath := flag.String("output", "faces.json", "Chemin pour sauvegarder les résultats")
    flag.Parse()

    if *imagePath == "" {
        log.Fatal("Veuillez spécifier une image à analyser")
    }

    // Créer le détecteur
    detector, err := face.NewDetector(*cascadePath, *modelPath)
    if err != nil {
        log.Fatalf("Erreur lors de l'initialisation du détecteur: %v", err)
    }

    // Détecter les visages
    faces, err := detector.DetectFaces(*imagePath)
    if err != nil {
        log.Fatalf("Erreur lors de la détection: %v", err)
    }

    // Sauvegarder les résultats
    result, err := json.MarshalIndent(faces, "", "  ")
    if err != nil {
        log.Fatalf("Erreur lors de la sérialisation: %v", err)
    }

    if err := os.WriteFile(*outputPath, result, 0644); err != nil {
        log.Fatalf("Erreur lors de la sauvegarde: %v", err)
    }

    log.Printf("Détection terminée. %d visages trouvés. Résultats sauvegardés dans %s", 
        len(faces), *outputPath)
}
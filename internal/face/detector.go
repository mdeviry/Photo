package face

import (
    "fmt"
    "image"
    _ "image/jpeg"
    "os"
    
    pigo "github.com/esimov/pigo/core"
    tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Detector gère la détection et la reconnaissance des visages
type Detector struct {
    classifier  *pigo.Pigo
    faceNet    *tf.SavedModel
    minSize     int
    scaleFactor float64
}

// NewDetector crée un nouveau détecteur
func NewDetector(cascadePath, modelPath string) (*Detector, error) {
    // Charger le classificateur Pigo
    cascadeFile, err := os.ReadFile(cascadePath)
    if err != nil {
        return nil, fmt.Errorf("impossible de lire le fichier cascade: %v", err)
    }
    
    p := pigo.NewPigo()
    classifier, err := p.Unpack(cascadeFile)
    if err != nil {
        return nil, fmt.Errorf("impossible de décompresser la cascade: %v", err)
    }

    // Charger le modèle FaceNet
    model, err := tf.LoadSavedModel(modelPath, []string{"serve"}, nil)
    if err != nil {
        return nil, fmt.Errorf("impossible de charger le modèle FaceNet: %v", err)
    }

    return &Detector{
        classifier:  classifier,
        faceNet:    model,
        minSize:    DefaultMinSize,
        scaleFactor: DefaultScaleFactor,
    }, nil
}

// DetectFaces détecte les visages dans une image
func (d *Detector) DetectFaces(filename string) ([]Face, error) {
    // Ouvrir et décoder l'image
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("impossible d'ouvrir l'image: %v", err)
    }
    defer file.Close()
    
    img, err := pigo.DecodeImage(file)
    if err != nil {
        return nil, fmt.Errorf("impossible de décoder l'image: %v", err)
    }

    // Convertir en niveaux de gris pour Pigo
    pixels := pigo.RgbToGrayscale(img)
    cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y

    params := pigo.CascadeParams{
        MinSize:     d.minSize,
        MaxSize:     min(cols, rows),
        ShiftFactor: 0.1,
        ScaleFactor: d.scaleFactor,
        ImageParams: pigo.ImageParams{
            Pixels: pixels,
            Rows:   rows,
            Cols:   cols,
            Dim:    cols,
        },
    }

    // Détecter les visages
    dets := d.classifier.RunCascade(params, 0)
    dets = d.classifier.ClusterDetections(dets, 0.2)

    // Convertir les détections en Face
    faces := make([]Face, 0, len(dets))
    for _, det := range dets {
        if det.Q < float32(DefaultConfidenceThreshold) {
            continue
        }

        face := Face{
            Rows:  rows,
            Cols:  cols,
            Score: det.Q,
            Area: Area{
                Name: "face",
                X:    float32(det.Col) / float32(cols),
                Y:    float32(det.Row) / float32(rows),
                Size: float32(det.Scale) / float32(cols),
            },
        }

        // Extraire l'embedding si possible
        if embeddings, err := d.getEmbeddings(img, &face); err == nil {
            face.Embeddings = embeddings
        }

        faces = append(faces, face)
    }

    return faces, nil
}

// getEmbeddings extrait les embeddings d'un visage
func (d *Detector) getEmbeddings(img image.Image, face *Face) (Embeddings, error) {
    // Découper et redimensionner le visage
    faceImg := cropAndResize(img, face.Area, DefaultImageSize)
    
    // Convertir l'image en tensor
    tensor, err := imageToTensor(faceImg)
    if err != nil {
        return nil, err
    }

    // Faire l'inférence avec FaceNet
    result, err := d.faceNet.Session.Run(
        map[tf.Output]*tf.Tensor{
            d.faceNet.Graph.Operation("input").Output(0): tensor,
        },
        []tf.Output{
            d.faceNet.Graph.Operation("embeddings").Output(0),
        },
        nil,
    )

    if err != nil {
        return nil, err
    }

    // Convertir le résultat en Embeddings
    embeddings := make(Embeddings, 1)
    embeddings[0] = make(Embedding, len(result[0].Value().([]float32)))
    copy(embeddings[0], result[0].Value().([]float32))

    return embeddings, nil
}

// cropAndResize découpe et redimensionne une portion d'image
func cropAndResize(img image.Image, area Area, size int) image.Image {
    // Implémenter le découpage et le redimensionnement
    // ... (code de traitement d'image)
    return img
}

// imageToTensor convertit une image en tensor pour TensorFlow
func imageToTensor(img image.Image) (*tf.Tensor, error) {
    // Implémenter la conversion
    // ... (code de conversion)
    return nil, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
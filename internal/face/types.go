package face

// Face représente un visage détecté avec ses caractéristiques
type Face struct {
    Rows       int        `json:"rows,omitempty"`
    Cols       int        `json:"cols,omitempty"`
    Score      float32    `json:"score,omitempty"`
    Area       Area       `json:"face,omitempty"`
    Embeddings Embeddings `json:"embeddings,omitempty"`
}

// Area représente une zone de l'image avec des coordonnées normalisées
type Area struct {
    Name  string  `json:"name,omitempty"`
    X     float32 `json:"x,omitempty"`
    Y     float32 `json:"y,omitempty"`
    Size  float32 `json:"size,omitempty"`
}

// Embedding représente un vecteur d'embedding facial
type Embedding []float32

// Embeddings est une collection d'embeddings
type Embeddings []Embedding

package face

// CropArea représente la zone de découpage relative pour un visage
type CropArea struct {
    Name   string
    X      float32
    Y      float32
    Width  float32
    Height float32
}

func NewCropArea(name string, x, y, width, height float32) CropArea {
    return CropArea{
        Name:   name,
        X:      x,
        Y:      y,
        Width:  width,
        Height: height,
    }
}
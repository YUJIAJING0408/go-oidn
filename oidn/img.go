package oidn

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func LoadPNG(path string) ([]float32, int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, 0, err
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	data := make([]float32, w*h*3)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			i := (y*w + x) * 3
			data[i+0] = float32(r) / 65535.0
			data[i+1] = float32(g) / 65535.0
			data[i+2] = float32(b) / 65535.0
		}
	}
	return data, w, h, nil
}

func SavePNG(data []float32, width, height int, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 3
			clamp := func(v float32) uint8 {
				if v < 0 {
					v = 0
				}
				if v > 1 {
					v = 1
				}
				return uint8(v * 255)
			}
			img.SetRGBA(x, y, color.RGBA{
				R: clamp(data[i+0]),
				G: clamp(data[i+1]),
				B: clamp(data[i+2]),
				A: 255,
			})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

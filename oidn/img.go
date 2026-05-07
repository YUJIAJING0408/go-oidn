package oidn

import (
	"fmt"
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
			img.SetRGBA(x, y, color.RGBA{
				R: clampAndScale(data[i+0]),
				G: clampAndScale(data[i+1]),
				B: clampAndScale(data[i+2]),
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

// clampAndScale 辅助函数（如果已在 img.go 中定义，无需重复）
func clampAndScale(v float32) uint8 {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return uint8(v*255.0 + 0.5) // 四舍五入可选项，原 clampAndScale 是直接截断，均可
}

// LoadPNGTilesWithOverlap 加载 PNG 并切分为 tilesX×tilesY 个均匀瓦片，
// 每个瓦片向外扩展 overlap 像素（超出图像边界的部分用黑色填充）。
// 返回：
//
//	tiles   - 瓦片列表，长度为 tilesX*tilesY
//	tileW   - 每个瓦片的宽度（baseW + 2*overlap）
//	tileH   - 每个瓦片的高度（baseH + 2*overlap）
//	origW   - 原图宽度
//	origH   - 原图高度
//
// 注意：原图宽高必须能被 tilesX 和 tilesY 整除。
func LoadPNGTilesWithOverlap(path string, tilesX, tilesY, overlap int) ([][]float32, int, int, int, int, error) {
	// 加载原图
	data, w, h, err := LoadPNG(path)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	if w%tilesX != 0 || h%tilesY != 0 {
		return nil, 0, 0, 0, 0, fmt.Errorf("image size %dx%d must be divisible by tilesX=%d, tilesY=%d", w, h, tilesX, tilesY)
	}

	baseW := w / tilesX
	baseH := h / tilesY
	tileW := baseW + 2*overlap
	tileH := baseH + 2*overlap
	numTiles := tilesX * tilesY
	tiles := make([][]float32, numTiles)

	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			idx := ty*tilesX + tx
			tile := make([]float32, tileW*tileH*3)

			// 当前瓦片在原始图像中的起始坐标（不含重叠）
			baseX := tx * baseW
			baseY := ty * baseH

			// 瓦片覆盖的原图区域（含重叠，可能超出图像）
			for y := 0; y < tileH; y++ {
				for x := 0; x < tileW; x++ {
					// 源坐标
					srcX := baseX - overlap + x
					srcY := baseY - overlap + y

					var r, g, b float32
					if srcX >= 0 && srcX < w && srcY >= 0 && srcY < h {
						// 在图像内部，读取实际像素
						srcIdx := (srcY*w + srcX) * 3
						r = data[srcIdx]
						g = data[srcIdx+1]
						b = data[srcIdx+2]
					} else {
						// 超出边界，填充黑色
						r, g, b = 0, 0, 0
					}

					dstIdx := (y*tileW + x) * 3
					tile[dstIdx] = r
					tile[dstIdx+1] = g
					tile[dstIdx+2] = b
				}
			}
			tiles[idx] = tile
		}
	}
	return tiles, tileW, tileH, w, h, nil
}

// SavePNGTilesWithOverlap 将带重叠的瓦片合并为原始分辨率图像并保存为 PNG。
// tiles: 瓦片列表（顺序同加载函数）
// tilesX, tilesY: 分块网格数
// overlap: 瓦片重叠像素数
// origW, origH: 原始图像的宽高
func SavePNGTilesWithOverlap(tiles [][]float32, tilesX, tilesY, overlap, origW, origH int, path string) error {
	if len(tiles) != tilesX*tilesY {
		return fmt.Errorf("tile count mismatch: got %d, expected %d", len(tiles), tilesX*tilesY)
	}
	if origW%tilesX != 0 || origH%tilesY != 0 {
		return fmt.Errorf("original size %dx%d must be divisible by tilesX=%d, tilesY=%d", origW, origH, tilesX, tilesY)
	}

	baseW := origW / tilesX
	baseH := origH / tilesY
	expectedTileW := baseW + 2*overlap
	expectedTileH := baseH + 2*overlap
	expectedLen := expectedTileW * expectedTileH * 3

	// 校验瓦片尺寸
	for i, tile := range tiles {
		if len(tile) != expectedLen {
			return fmt.Errorf("tile %d size mismatch: got %d, expected %d (tileW=%d, tileH=%d)", i, len(tile), expectedLen, expectedTileW, expectedTileH)
		}
	}

	// 创建输出图像
	img := image.NewRGBA(image.Rect(0, 0, origW, origH))

	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			tile := tiles[ty*tilesX+tx]
			// 有效区域在瓦片内的起始偏移（跳过重叠部分）
			startX := overlap
			startY := overlap
			destX := tx * baseW
			destY := ty * baseH

			for y := 0; y < baseH; y++ {
				for x := 0; x < baseW; x++ {
					srcIdx := ((startY+y)*expectedTileW + (startX + x)) * 3
					r := clampAndScale(tile[srcIdx])
					g := clampAndScale(tile[srcIdx+1])
					b := clampAndScale(tile[srcIdx+2])
					img.SetRGBA(destX+x, destY+y, color.RGBA{R: r, G: g, B: b, A: 255})
				}
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

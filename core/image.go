package core

import (
	"image"
	"image/draw"
	"image/gif"
	"os"

	"github.com/nfnt/resize"
)

func convertAndResizeAndCenter(w, h int, src image.Image, xOffsetPercent float64) *image.Gray {
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(image.Point{int(float64(w*2)*xOffsetPercent) - w, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func Gif2Animation(w int, h int, path string) ([]*image.Gray, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}

	animationStep := 0.05
	animation := make([]*image.Gray, int(1/animationStep))
	currentFrame := 0
	for i := 0; i < len(animation); i++ {
		animation[i] = convertAndResizeAndCenter(w, h, g.Image[currentFrame], float64(i)*animationStep)
		currentFrame = (currentFrame + 1) % len(g.Image)
	}

	return animation, nil
}

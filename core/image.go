package core

import (
	"image"
	"image/draw"
	"image/gif"
	"os"
	"time"
)

type Animation struct {
	Frames        []*image.Gray
	FrameDuration []time.Duration
	Sequence      []*image.Point
}

func convertToGray(src image.Image) *image.Gray {
	r := src.Bounds()
	img := image.NewGray(r)
	r = r.Add(image.Point{})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func CreateFrame(w int, h int, src *image.Gray, offset image.Point) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(offset)
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func Gif2Animation(w int, h int, path string, duration time.Duration) (*Animation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}

	animation := Animation{
		Frames:        make([]*image.Gray, len(g.Image)),
		FrameDuration: make([]time.Duration, len(g.Image)),
		Sequence:      []*image.Point{},
	}

	totalAnimationDuration := time.Duration(0)
	for currentFrame := 0; currentFrame < len(g.Image); currentFrame++ {
		animation.Frames[currentFrame] = convertToGray(g.Image[currentFrame])
		animation.FrameDuration[currentFrame] = time.Duration(10*g.Delay[currentFrame]) * time.Millisecond
		totalAnimationDuration += animation.FrameDuration[currentFrame]
	}

	totalWidth := g.Config.Width + w
	totalLoops := float64(duration / totalAnimationDuration)

	increment := totalWidth / (int(totalLoops) * len(g.Image))

	for i := -g.Config.Width; i <= w; i += increment {
		animation.Sequence = append(animation.Sequence, &image.Point{i, 0})
	}

	return &animation, nil
}

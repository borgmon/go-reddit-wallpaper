package main

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/disintegration/imaging"
)

const (
	threshold     = .45
	dimLevel      = -10
	compressLevel = 90
)

// calculate brightness using RMS of grayscaled picture
func checkDarkImage(img image.Image) (bool, error) {
	newImage := imaging.Grayscale(img)
	var sum int
	for i, v := range newImage.Pix {
		if i%4 == 0 {
			sum += int(v) * int(v)
		}
	}
	a := sum / (newImage.Rect.Max.X * newImage.Rect.Max.Y)

	result := math.Sqrt(float64(a))

	darkScale := result / float64(255)

	return darkScale < threshold, nil
}

func dimImage(img image.Image) image.Image {
	return imaging.AdjustBrightness(img, dimLevel)
}

func getDimensions(img image.Image) (int, int, error) {
	return img.Bounds().Dx(), img.Bounds().Dy(), nil
}

func encodePNG(img image.Image) ([]byte, error) {
	encodeOption := imaging.JPEGQuality(compressLevel)

	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.JPEG, encodeOption)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), nil
}

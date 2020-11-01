package main

import (
	"bytes"
	"log"
	"math"

	"github.com/disintegration/imaging"
)

const (
	threshold = .45
	dimLevel  = -20
)

// calculate brightness using RMS of grayscaled picture
func checkDarkImage(img []byte) (bool, error) {
	decodedImage, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return false, err
	}
	newImage := imaging.Grayscale(decodedImage)
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

func DimImage(img []byte) ([]byte, error) {
	decodedImage, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}
	newImage := imaging.AdjustBrightness(decodedImage, dimLevel)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, newImage, imaging.PNG)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), nil
}

func getDimensions(img []byte) (int, int, error) {
	decodedImage, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return 0, 0, err
	}
	return decodedImage.Bounds().Dx(), decodedImage.Bounds().Dy(), nil
}

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

// func main() {
// 	// res, err := http.Get("https://coolbackgrounds.io/images/backgrounds/white/pure-white-background-85a2a7fd.jpg")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer res.Body.Close()
// 	// d, err := ioutil.ReadAll(res.Body)

// 	src, err := imaging.Open("w.png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	newImage := imaging.Grayscale(src)
// 	// 0 0 0 0
// 	var sum int
// 	for i, v := range newImage.Pix {
// 		if (i+1)%4 == 0 {
// 			sum += int(v) * int(v)
// 		}
// 	}
// 	a := sum / (newImage.Rect.Max.X * newImage.Rect.Max.Y)

// 	result := math.Sqrt(float64(a))
// 	fmt.Println(result)
// }

func CheckDarkImage(img []byte) (bool, error) {
	decodedImage, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return false, err
	}
	newImage := imaging.Grayscale(decodedImage)
	var sum int
	for i, v := range newImage.Pix {
		if (i+1)%4 == 0 {
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

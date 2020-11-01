package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/kkyr/wallpaper"
)

const (
	randomMax  = 20
	fetchLimit = 20
)

type savedImage struct {
	url     string
	byteArr []byte
	isDark  bool
}

func Start() {
	savedSubreddit := mainApp.Preferences().String("subreddits")
	savedSubreddit = trimWhiteSpace(savedSubreddit)
	savedSorting := mainApp.Preferences().String("sorting")
	savedPreferDarker := mainApp.Preferences().String("prefer_darker")
	savedDeepscan := mainApp.Preferences().Bool("deepscan")

	subreddit := randomElement(strings.Split(savedSubreddit, ","))

	randomIndex, deathCounter := 0, 0
	if mainApp.Preferences().String("first_or_random") == "random" {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(randomMax-1) + 1
	}

	var finalImage, lastImage *savedImage = nil, nil
	afterID := ""

	for finalImage == nil {
		newLogInfo("Getting a new page from: " + subreddit + " , afterID=" + afterID)
		payload, newAfterID, err := getReddit(subreddit, savedSorting, afterID)
		if err != nil {
			newLogError("cannot get reddit", err)
			return
		}
		afterID = newAfterID
		for _, v := range payload.Data.Children {
			var image *savedImage = nil

			if v.Data.Preview.Images == nil {
				if savedDeepscan {
					image, err = getImageDeepscan(&v)
					if err != nil {
						newLogError("ImageDeepscan error", err)
						continue
					}
				} else {
					continue
				}
			} else {
				image, err = getImage(&v)
				if err != nil {
					newLogError("get image error", err)
					continue
				}
			}

			if image != nil {

				if savedPreferDarker == "only dark images" || savedPreferDarker == "dim images" {
					isDark, err := checkDarkImage(image.byteArr)
					if err != nil {
						newLogError("check dark image error", err)
						continue
					}
					image.isDark = isDark
					if !isDark && savedPreferDarker != "dim images" {
						newLogInfo("Skip image because it's not dark: " + image.url)
						continue
					}
				}

				lastImage = image
				randomIndex--
				if randomIndex <= 0 {
					finalImage = lastImage
					break
				}
			}
		}
		deathCounter++
		if deathCounter > fetchLimit {
			break
		}
	}

	if finalImage == nil {
		if lastImage != nil {
			finalImage = lastImage
		} else {
			newLogError("", errors.New("No image found"))
			return
		}
	}

	if savedPreferDarker == "dim images" {
		newLogInfo("Dimming brightness")
		byteArr, err := dimImage(finalImage.byteArr)
		if err != nil {
			newLogError("check dark image error", err)
			return
		}
		finalImage.byteArr = byteArr
	}

	path, err := saveWallperFromByteArray(finalImage.byteArr)
	if err != nil {
		newLogError("set wallpaper error", err)
	} else {
		newLogInfo("Successfully set wallpaper: " + finalImage.url)
		newLogInfo("Saved to: " + path)
	}
}

func randomElement(ls []string) string {
	randomIndex := rand.Intn(len(ls))
	return ls[randomIndex]
}

func trimWhiteSpace(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)
}

func saveWallperFromByteArray(byteArr []byte) (string, error) {
	tmpfile, err := ioutil.TempFile("", "go-reddit-wallpaper-temp-")
	if err != nil {
		return "", err
	}

	_, err = tmpfile.Write(byteArr)
	if err != nil {
		return "", err
	}
	err = tmpfile.Close()
	if err != nil {
		return "", err
	}
	err = wallpaper.SetFromFile(tmpfile.Name())
	return tmpfile.Name(), err
}

func getImage(v *PayloadDataChild) (*savedImage, error) {
	minWidth := mainApp.Preferences().Int("min_width")
	minHeight := mainApp.Preferences().Int("min_height")
	width := v.Data.Preview.Images[0].Source.Width
	height := v.Data.Preview.Images[0].Source.Height
	if width >= minWidth && height >= minHeight {
		url := fixPreviewUrl(v.Data.Preview.Images[0].Source.Url)
		byteArr, err := download(url)
		if err != nil {
			return nil, err
		}
		return &savedImage{url: url, byteArr: byteArr}, nil
	}
	return nil, nil
}

func getImageDeepscan(v *PayloadDataChild) (*savedImage, error) {
	minWidth := mainApp.Preferences().Int("min_width")
	minHeight := mainApp.Preferences().Int("min_height")

	url := v.Data.Url

	if url == "" {
		return nil, nil
	}

	if url[len(url)-4:] != ".png" && url[len(url)-4:] != ".jpg" && url[len(url)-5:] != ".jpeg" {
		return nil, nil
	}

	byteArr, err := download(url)
	if err != nil {
		return nil, err
	}

	width, height, err := getDimensions(byteArr)
	if err != nil {
		return nil, err
	}
	if width >= minWidth && height >= minHeight {
		return &savedImage{url: fixPreviewUrl(url), byteArr: byteArr}, nil
	}
	return nil, nil
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func fixPreviewUrl(url string) string {
	return strings.Replace(url, "amp;", "", -1)
}

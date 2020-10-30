package main

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/kkyr/wallpaper"
)

const (
	randomMax = 20
)

func Start() {
	savedSubreddit := MainApp.Preferences().String("subreddits")
	savedSubreddit = trimWhiteSpace(savedSubreddit)
	savedSorting := MainApp.Preferences().String("sorting")

	subreddit := randomElement(strings.Split(savedSubreddit, ","))

	var body struct {
		image string
	}

	randomIndex := 1
	if MainApp.Preferences().String("first_or_random") == "random" {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(randomMax-1) + 1
	}

	for body.image == "" {
		res, err := GetReddit(subreddit, savedSorting)
		if err != nil {
			NewLogError(err)
		}
		NewLogInfo("Getting a new page from: " + subreddit + " , afterID=" + AfterId)
		image := ""
		lastImage := ""
		for _, v := range res.Data.Children {
			if v.Data.Preview.Images == nil {
				continue
			}

			result := getImage(&v)

			if result != "" {
				lastImage = result
			}
			randomIndex--
			if randomIndex <= 0 {
				image = lastImage
				break
			}
		}

		if image == "" {
			continue
		}

		if err != nil {
			NewLogError(err)
		}
		if image != "" {
			body.image = image
			AfterId = ""
		}

	}

	err := wallpaper.SetFromURL(body.image)
	if err != nil {
		NewLogError(err)
	} else {
		NewLogInfo("Successfully set wallpaper: " + body.image)
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

func getImage(v *PayloadDataChild) string {
	minWidth := MainApp.Preferences().Int("min_width")
	minHeight := MainApp.Preferences().Int("min_height")
	width := v.Data.Preview.Images[0].Source.Width
	height := v.Data.Preview.Images[0].Source.Height
	if width >= minWidth && height >= minHeight {
		return fixPreviewUrl(v.Data.Preview.Images[0].Source.Url)
	}
	return ""
}

func fixPreviewUrl(url string) string {
	return strings.Replace(url, "amp;", "", -1)
}

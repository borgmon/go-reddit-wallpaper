package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/reujab/wallpaper"
)

func Start() {
	savedSubreddit := MainApp.Preferences().String("subreddits")
	savedSubreddit = trimWhiteSpace(savedSubreddit)
	savedSorting := MainApp.Preferences().String("sorting")

	subreddit := randomElement(strings.Split(savedSubreddit, ","))

	var body struct {
		payload *RedditPayload
		image   string
	}
	for body.image == "" {
		res, err := GetReddit(subreddit, savedSorting)
		if err != nil {
			ErrorPopup(err)
		}
		fmt.Println(res.Data.Children[0].Data.Name)
		image, err := getImage(res)
		if err != nil {
			ErrorPopup(err)
		}
		if image != "" {
			body.image = image
			AfterId = ""
		}

	}

	err := wallpaper.SetFromURL(body.image)
	if err != nil {
		ErrorPopup(err)
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

func getImage(payload *RedditPayload) (string, error) {
	for _, v := range payload.Data.Children {
		minWidthStr := MainApp.Preferences().String("min_width")
		minWidth, err := strconv.Atoi(minWidthStr)
		if err != nil {
			ErrorPopup(err)
		}
		minHeightStr := MainApp.Preferences().String("min_height")
		minHeight, err := strconv.Atoi(minHeightStr)
		if err != nil {
			ErrorPopup(err)
		}
		width := v.Data.Preview.Images[0].Source.Width
		height := v.Data.Preview.Images[0].Source.Height
		if width >= minWidth && height >= minHeight {
			return v.Data.Url, nil
		}
	}
	return "", errors.New("No images met requirement")
}
package main

import (
	"encoding/json"
	"net/http"
)

type RedditPayload struct {
	Data PayloadData
}

type PayloadData struct {
	Children []PayloadDataChild
}

type PayloadDataChild struct {
	Data PayloadDataChildData
}

type PayloadDataChildData struct {
	Preview PayloadDataChildDataPreview
	Name    string
}

type PayloadDataChildDataPreview struct {
	Images []PayloadDataChildDataPreviewImage
}
type PayloadDataChildDataPreviewImage struct {
	Source PayloadDataChildDataPreviewImageSource
}

type PayloadDataChildDataPreviewImageSource struct {
	Url    string
	Width  int
	Height int
}

var AfterId = ""

func GetReddit(subreddit, sort string) (result *RedditPayload, err error) {
	client := &http.Client{}
	var url string
	if AfterId == "" {
		url = "https://www.reddit.com/" + subreddit + "/" + sort + ".json?count=25"
	} else {
		url = "https://www.reddit.com/" + subreddit + "/" + sort + ".json?count=25&after=" + AfterId
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-reddit-wallpaper/1.0")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	AfterId = result.Data.Children[0].Data.Name
	return
}

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
	Url     string
}

type PayloadDataChildDataPreview struct {
	Images []PayloadDataChildDataPreviewImage
}
type PayloadDataChildDataPreviewImage struct {
	Source PayloadDataChildDataPreviewImageSource
}

type PayloadDataChildDataPreviewImageSource struct {
	Width  int16
	Height int16
}

func GetReddit(subreddit, sort string) (result *RedditPayload, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.reddit.com/"+subreddit+"/"+sort+".json", nil)
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
	return
}

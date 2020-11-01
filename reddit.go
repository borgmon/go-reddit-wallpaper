package main

import (
	"encoding/json"
	"errors"
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
	Url     string
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

func getReddit(subreddit, sort string, afterID string) (result *RedditPayload, newAfterID string, err error) {
	client := &http.Client{}
	var url string
	if afterID == "" {
		url = "https://www.reddit.com/" + subreddit + "/" + sort + ".json?count=25"
	} else {
		url = "https://www.reddit.com/" + subreddit + "/" + sort + ".json?count=25&after=" + afterID
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "go-reddit-wallpaper/"+version)
	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, "", err
	}
	if len(result.Data.Children) == 0 {
		return nil, "", errors.New("reach end of page limit")
	}
	newAfterID = result.Data.Children[0].Data.Name
	return
}

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/reujab/wallpaper"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/wallpapers/top.json", nil)
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	var result RedditPayload
	json.NewDecoder(res.Body).Decode(&result)
	err = wallpaper.SetFromURL(result.Data.Children[0].Data.Url)
	if err != nil {
		log.Fatalln(err)
	}
}

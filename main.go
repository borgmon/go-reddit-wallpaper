package main

import (
	"log"

	"github.com/reujab/wallpaper"
)

func main() {
	res, err := GetReddit("r/wallpapers", "top")
	if err != nil {
		log.Fatalln(err)
	}

	err = wallpaper.SetFromURL(res.Data.Children[0].Data.Url)
	if err != nil {
		log.Fatalln(err)
	}
}

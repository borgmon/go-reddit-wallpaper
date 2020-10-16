package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	settingWindow := app.NewWindow("Preference")

	settingWindow.SetContent(widget.NewLabel("Subreddits"))

}

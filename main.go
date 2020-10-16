package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

const (
	buildInSubreddits = "r/wallpaper,r/wallpapers"
)

var (
	goApp         = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	settingWindow = goApp.NewWindow("Preference")
	sorting       = []string{"top", "hot", "new"}
)

func main() {

	settingWindow.SetFixedSize(true)
	settingWindow.CenterOnScreen()

	subredditsEntry := getInputBox("subreddits", buildInSubreddits)

	minWidthEntry := getInputBox("min_width", "1920")
	minHeightEntry := getInputBox("min_height", "1080")
	minSizeBox := widget.NewHBox(minWidthEntry, widget.NewLabel("x"), minHeightEntry)

	intervalEntry := getInputBox("interval", "@daily")

	sortingSelect := widget.NewSelect(sorting, func(text string) {
		goApp.Preferences().SetString("sorting", text)
	})
	sortingSelect.SetSelected(goApp.Preferences().StringWithFallback("sorting", sorting[0]))

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		goApp.Preferences().SetBool("autorun", toggle)
	})
	autorunCheck.SetChecked(goApp.Preferences().BoolWithFallback("autorun", true))

	settingWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widget.NewLabelWithStyle("Subreddits", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subredditsEntry,
		widget.NewLabelWithStyle("Minimum Size", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		minSizeBox,
		widget.NewLabelWithStyle("Refresh Interval", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		intervalEntry,
		widget.NewLabelWithStyle("Sorting Method", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sortingSelect,
		widget.NewLabelWithStyle("Auto Run", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		autorunCheck,
	))
	settingWindow.ShowAndRun()
}

func getInputBox(name, fallback string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(goApp.Preferences().StringWithFallback(name, fallback))
	entry.SetPlaceHolder(fallback)
	entry.OnChanged = func(text string) {
		goApp.Preferences().SetString(name, text)
	}
	return entry
}

package main

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var (
	MainApp           = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	SettingWindow     = MainApp.NewWindow("Preference")
	sorting           = []string{"top", "hot", "new"}
	buildInSubreddits = "r/wallpaper,r/wallpapers"
)

func main() {
	SettingWindow.SetFixedSize(true)
	SettingWindow.CenterOnScreen()

	subredditsEntry := getStringInputBox("subreddits", buildInSubreddits)

	minWidthEntry := getStringInputBox("min_width", "1920")
	minHeightEntry := getStringInputBox("min_height", "1080")
	minSizeBox := widget.NewHBox(minWidthEntry, widget.NewLabel("x"), minHeightEntry)

	intervalEntry := getStringInputBox("interval", "@daily")

	sortingSelect := widget.NewSelect(sorting, func(text string) {
		MainApp.Preferences().SetString("sorting", text)
	})
	sortingSelect.SetSelected(MainApp.Preferences().StringWithFallback("sorting", sorting[0]))

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		MainApp.Preferences().SetBool("autorun", toggle)
		// TODO: set cron job
	})
	autorunCheck.SetChecked(MainApp.Preferences().BoolWithFallback("autorun", true))

	SettingWindow.SetContent(widget.NewVBox(
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
		widget.NewButton("run", func() {
			Start()
		}),
	))

	SettingWindow.ShowAndRun()

}

func ErrorPopup(err error) {
	w := MainApp.NewWindow("Error")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(300, 200))
	w.SetContent(widget.NewScrollContainer(widget.NewLabel(err.Error())))
	w.Show()

}

func getStringInputBox(name, fallback string) *widget.Entry {
	entry := widget.NewEntry()
	value := MainApp.Preferences().StringWithFallback(name, fallback)
	MainApp.Preferences().SetString(name, value)
	entry.SetText(value)
	entry.SetPlaceHolder(fallback)
	entry.OnChanged = func(text string) {
		MainApp.Preferences().SetString(name, text)
	}
	return entry
}

func getIntInputBox(name string, fallback int) *widget.Entry {
	entry := widget.NewEntry()
	value := MainApp.Preferences().IntWithFallback(name, fallback)
	MainApp.Preferences().SetInt(name, value)
	text := strconv.Itoa(value)
	entry.SetText(text)
	entry.SetPlaceHolder(text)
	entry.OnChanged = func(text string) {
		i, _ := strconv.Atoi(text)
		MainApp.Preferences().SetInt(name, i)
	}
	return entry
}

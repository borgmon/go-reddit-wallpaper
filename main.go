package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
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

	subredditsEntry := getInputBox("subreddits", buildInSubreddits)

	minWidthEntry := getInputBox("min_width", "1920")
	minHeightEntry := getInputBox("min_height", "1080")
	minSizeBox := widget.NewHBox(minWidthEntry, widget.NewLabel("x"), minHeightEntry)

	intervalEntry := getInputBox("interval", "@daily")

	sortingSelect := widget.NewSelect(sorting, func(text string) {
		MainApp.Preferences().SetString("sorting", text)
	})
	sortingSelect.SetSelected(MainApp.Preferences().StringWithFallback("sorting", sorting[0]))

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		MainApp.Preferences().SetBool("autorun", toggle)
	})
	autorunCheck.SetChecked(MainApp.Preferences().BoolWithFallback("autorun", true))

	SettingWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
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

func getInputBox(name, fallback string) *widget.Entry {
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

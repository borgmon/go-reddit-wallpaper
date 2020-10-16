package main

import (
	"errors"
	"runtime"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/ProtonMail/go-autostart"
)

var (
	MainApp           = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	SettingWindow     = MainApp.NewWindow("Preference")
	sorting           = []string{"top", "hot", "new"}
	firstOrRandom     = []string{"first", "random"}
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

	firstOrRandomSelect := widget.NewSelect(firstOrRandom, func(text string) {
		MainApp.Preferences().SetString("first_or_random", text)
	})
	firstOrRandomSelect.SetSelected(MainApp.Preferences().StringWithFallback("first_or_random", firstOrRandom[0]))

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		MainApp.Preferences().SetBool("autorun", toggle)
		_, file, _, ok := runtime.Caller(1)
		if !ok {
			ErrorPopup(errors.New("Autorun setup failed"))
		}
		autoStartApp := &autostart.App{
			Name:        "go-reddit-wallpaper",
			DisplayName: "Go Reddit WallPaper",
			Exec:        []string{"bash", "-c", file + " >> ~/autostart.txt"},
		}
		if toggle {
			autoStartApp.Enable()
		} else {
			autoStartApp.Disable()
		}
	})
	autorunCheck.SetChecked(MainApp.Preferences().BoolWithFallback("autorun", false))

	SettingWindow.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Subreddits", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subredditsEntry,
		widget.NewLabelWithStyle("Minimum Size", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		minSizeBox,
		widget.NewLabelWithStyle("Refresh Interval", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		intervalEntry,
		widget.NewLabelWithStyle("Sorting Method", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sortingSelect,
		widget.NewLabelWithStyle("Select First Or Random", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		firstOrRandomSelect,
		widget.NewLabelWithStyle("Auto Run (experimental)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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
	w.SetOnClosed(func() {
		MainApp.Quit()
	})
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
		i, err := strconv.Atoi(text)
		if err != nil {
			ErrorPopup(err)
		}
		MainApp.Preferences().SetInt(name, i)
	}
	return entry
}

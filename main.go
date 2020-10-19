package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/ProtonMail/go-autostart"
	"github.com/getlantern/systray"
	"github.com/robfig/cron/v3"
)

var (
	MainApp           = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	sorting           = []string{"top", "hot", "new"}
	firstOrRandom     = []string{"first", "random"}
	buildInSubreddits = "r/wallpaper,r/wallpapers"
	iconPath          = "./Icon.png"
	IconRecource      fyne.Resource
	prefWindowChannel = make(chan bool)
	settingWindow     fyne.Window
	cronJob           = cron.New()
)

func main() {
	cronJob.Start()
	SetupGUI()
	go systray.Run(onReady, onExit)
	settingWindow = BuildPrefWindow()
	go Start()
	MainApp.Run()
}

func SetupGUI() {
	iconRecource, err := fyne.LoadResourceFromPath(iconPath)
	IconRecource = iconRecource
	if err != nil {
		log.Fatalln(err)
	}
	MainApp.SetIcon(IconRecource)
}

func BuildPrefWindow() fyne.Window {
	settingWindow := MainApp.NewWindow("Preferences")
	settingWindow.SetIcon(IconRecource)
	settingWindow.SetFixedSize(true)
	settingWindow.CenterOnScreen()
	settingWindow.SetCloseIntercept(func() {
		settingWindow.Hide()
	})

	subredditsEntry := getStringInputBox("subreddits", buildInSubreddits)

	minSizeErrorLabel := widget.NewLabel("")

	minWidthEntry := getIntInputBox("min_width", 1920, minSizeErrorLabel)
	minHeightEntry := getIntInputBox("min_height", 1080, minSizeErrorLabel)
	minSizeBox := widget.NewHBox(minWidthEntry, widget.NewLabel("x"), minHeightEntry)

	intervalEntryErrorLabel := widget.NewLabel("")
	intervalEntry := widget.NewEntry()
	value := MainApp.Preferences().StringWithFallback("interval", "@daily")
	MainApp.Preferences().SetString("interval", value)
	intervalEntry.SetText(value)
	intervalEntry.SetPlaceHolder("@daily")

	intervalEntry.OnChanged = func(text string) {
		_, err := cron.ParseStandard(text)
		if err != nil {
			intervalEntryErrorLabel.SetText("Wrong Format")
		} else {
			intervalEntryErrorLabel.SetText("")
			MainApp.Preferences().SetString("interval", text)
			clearAllCronJobs()
			cronJob.AddFunc(text, func() {
				go Start()
			})
		}
	}

	sortingSelect := widget.NewSelect(sorting, func(text string) {
		MainApp.Preferences().SetString("sorting", text)
	})
	sortingSelect.SetSelected(MainApp.Preferences().StringWithFallback("sorting", sorting[0]))
	sortingSelect.OnChanged = func(text string) {
		MainApp.Preferences().SetString("sorting", text)
		go Start()
	}

	firstOrRandomSelect := widget.NewSelect(firstOrRandom, func(text string) {
		MainApp.Preferences().SetString("first_or_random", text)
	})
	firstOrRandomSelect.SetSelected(MainApp.Preferences().StringWithFallback("first_or_random", firstOrRandom[0]))
	firstOrRandomSelect.OnChanged = func(text string) {
		MainApp.Preferences().SetString("first_or_random", text)
		go Start()
	}

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

	settingWindow.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Subreddits", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subredditsEntry,
		widget.NewLabelWithStyle("Minimum Size", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		minSizeBox,
		minSizeErrorLabel,
		widget.NewLabelWithStyle("Refresh Interval", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		intervalEntry,
		intervalEntryErrorLabel,
		widget.NewLabelWithStyle("Sorting Method", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sortingSelect,
		widget.NewLabelWithStyle("Select First Or Random", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		firstOrRandomSelect,
		widget.NewLabelWithStyle("Auto Run (experimental)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		autorunCheck,
	))

	settingWindow.Show()
	return settingWindow
}

func ErrorPopup(err error) {
	w := MainApp.NewWindow("Error")
	w.SetIcon(IconRecource)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(300, 200))
	w.SetContent(widget.NewScrollContainer(widget.NewLabel(err.Error())))
	w.Show()
	// w.SetOnClosed(func() {
	// 	MainApp.Quit()
	// })
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

func getIntInputBox(name string, fallback int, errorMsg *widget.Label) *widget.Entry {
	entry := widget.NewEntry()
	value := MainApp.Preferences().IntWithFallback(name, fallback)
	MainApp.Preferences().SetInt(name, value)
	text := strconv.Itoa(value)
	entry.SetText(text)
	entry.SetPlaceHolder(text)
	entry.OnChanged = func(text string) {
		i, err := strconv.Atoi(text)
		if err != nil {
			errorMsg.SetText("Not a number")
		} else {
			MainApp.Preferences().SetInt(name, i)
			errorMsg.SetText("")
		}
	}
	return entry
}

func onReady() {
	systray.SetIcon(IconRecource.Content())
	systray.SetTitle("Go Reddit WallPaper")
	systray.SetTooltip("Go Reddit WallPaper")

	mQuit := systray.AddMenuItem("Quit", "Quit Go Reddit WallPaper")
	mPref := systray.AddMenuItem("Preferences", "Change Preferences")
	mRefresh := systray.AddMenuItem("Refresh Now", "Refresh Now!")
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			return

		case <-mPref.ClickedCh:
			fmt.Println("here")
			settingWindow.Show()

		case <-mRefresh.ClickedCh:
			go Start()
		}

	}
}
func onExit() {
	MainApp.Quit()
}

func clearAllCronJobs() {
	jobs := cronJob.Entries()
	for _, job := range jobs {
		cronJob.Remove(job.ID)
	}
}

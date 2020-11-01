package main

import (
	"net/url"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

const (
	githubLink = "https://github.com/borgmon/go-reddit-wallpaper"
	version    = "v1.3"
)

var (
	MainApp           = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	sorting           = []string{"top", "hot", "new"}
	firstOrRandom     = []string{"first", "random"}
	preferDarker      = []string{"none", "only dark images"}
	buildInSubreddits = "r/wallpaper,r/wallpapers"
	logWindow         fyne.Window
	LogEntry          *widget.Entry
	settingWindow     fyne.Window
	cronJob           = newCron()
	trayIconResource  []byte
)

func main() {
	cronJob.Start()
	SetupIcon()
	go startTray()
	logWindow = BuildLogWindow()
	settingWindow = BuildPrefWindow()
	go Start()
	MainApp.Run()
}

func SetupIcon() {
	MainApp.SetIcon(PngIconResource)
	// windows tray icon issue walk around https://github.com/reujab/wallpaper/pull/15
	if runtime.GOOS == "windows" {
		trayIconResource = IcoIconResource.StaticContent
	} else {
		trayIconResource = PngIconResource.StaticContent
	}
}

func BuildPrefWindow() fyne.Window {
	settingWindow := MainApp.NewWindow("Preferences")
	settingWindow.SetIcon(PngIconResource)
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
	url, err := url.Parse("https://godoc.org/github.com/robfig/cron")
	checkError("parse cron doc error", err)
	intervalLink := widget.NewHyperlink("See example", url)
	value := MainApp.Preferences().StringWithFallback("interval", "@daily")
	MainApp.Preferences().SetString("interval", value)
	intervalEntry.SetText(value)
	intervalEntry.SetPlaceHolder("@daily")

	intervalEntry.OnChanged = func(text string) {
		_, err := parseCron(text)
		if err != nil {
			intervalEntryErrorLabel.SetText("Wrong Format")
		} else {
			intervalEntryErrorLabel.SetText("")
			MainApp.Preferences().SetString("interval", text)
			_, err := clearAndSetCron(text)
			checkError("set cron failed", err)
		}
	}

	sortingSelect := getSelect("sorting", sorting)

	firstOrRandomSelect := getSelect("first_or_random", firstOrRandom)

	perferDarkerSelect := getSelect("prefer_darker", preferDarker)

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		MainApp.Preferences().SetBool("autorun", toggle)

		autoStartApp, err := newAutoRun()
		checkError("get autostart app failed", err)
		if toggle {
			autoStartApp.Enable()
		} else {
			autoStartApp.Disable()
		}

	})
	autorunEnabled := MainApp.Preferences().BoolWithFallback("autorun", false)
	autorunCheck.SetChecked(autorunEnabled)

	deepscanCheck := widget.NewCheck("deepscan", func(toggle bool) {
		MainApp.Preferences().SetBool("deepscan", toggle)
	})
	deepscanEnabled := MainApp.Preferences().BoolWithFallback("deepscan", false)
	deepscanCheck.SetChecked(deepscanEnabled)

	settingWindow.SetContent(container.NewAdaptiveGrid(2,
		widget.NewVBox(widget.NewLabel("Subreddits")),
		widget.NewVBox(
			subredditsEntry,
		),

		widget.NewVBox(widget.NewLabel("Minimum Size")),
		widget.NewVBox(
			minSizeBox,
			minSizeErrorLabel,
		),

		widget.NewVBox(
			widget.NewLabel("Refresh Interval"),
			intervalLink,
		),
		widget.NewVBox(
			intervalEntry,
			intervalEntryErrorLabel,
		),

		widget.NewVBox(widget.NewLabel("Sorting Method")),
		widget.NewVBox(
			sortingSelect,
		),

		widget.NewVBox(widget.NewLabel("Select First Or Random")),
		widget.NewVBox(
			firstOrRandomSelect,
		),

		widget.NewVBox(widget.NewLabel("Prefer Darker")),
		widget.NewVBox(
			perferDarkerSelect,
		),

		widget.NewVBox(widget.NewLabel("Auto Run")),
		widget.NewVBox(
			autorunCheck,
		),

		widget.NewVBox(
			widget.NewLabel("Deep Scan"),
			widget.NewLabel("(download picture to check dimensions)"),
		),
		widget.NewVBox(
			deepscanCheck,
		),

		widget.NewVBox(widget.NewLabel("version: "+version)),
		widget.NewVBox(
			widget.NewButtonWithIcon("Github", GithubPngResource, func() {
				url, _ := url.Parse("https://github.com/borgmon/go-reddit-wallpaper")
				err = fyne.CurrentApp().OpenURL(url)
				checkError("open github url failed", err)
			}),
		),
	))

	if autorunEnabled {
		settingWindow.Hide()
	} else {
		settingWindow.Show()
	}

	return settingWindow
}
func BuildLogWindow() fyne.Window {
	logWindow := MainApp.NewWindow("Logs")
	logWindow.SetIcon(PngIconResource)
	logWindow.CenterOnScreen()
	logWindow.Resize(fyne.NewSize(600, 800))
	logWindow.SetCloseIntercept(func() {
		logWindow.Hide()
	})

	LogEntry = widget.NewMultiLineEntry()
	LogEntry.Disable()

	logWindow.SetContent(container.NewScroll(LogEntry))

	logWindow.Hide()
	return logWindow
}
func NewLogError(text string, err error) {
	LogEntry.Text += time.Now().Format(time.RFC3339) + "\tERROR\t" + text + "\t" + err.Error() + "\n"
}
func NewLogInfo(text string) {
	LogEntry.Text += time.Now().Format(time.RFC3339) + "\tINFO\t" + text + "\n"
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

func getSelect(name string, selection []string) *widget.Select {
	selectEl := widget.NewSelect(selection, func(text string) {
		MainApp.Preferences().SetString(name, text)
	})
	selectEl.SetSelected(MainApp.Preferences().StringWithFallback(name, selection[0]))
	selectEl.OnChanged = func(text string) {
		MainApp.Preferences().SetString(name, text)
		go Start()
	}
	return selectEl
}

func checkError(text string, err error) {
	if err != nil {
		NewLogError(text, err)
	}
}

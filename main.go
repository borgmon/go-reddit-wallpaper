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
	mainApp           = app.NewWithID("com.github.borgmon.go-reddit-wallpaper")
	sorting           = &options{"top": "Top", "hot": "Hot", "new": "New"}
	firstOrRandom     = &options{"first": "First", "random": "Random"}
	preferDarker      = &options{"none": "None", "dark_image": "Only use dark images", "dim_image": "Dim images (post processing)"}
	buildInSubreddits = "r/wallpaper,r/wallpapers"
	logWindow         fyne.Window
	logEntry          *widget.Entry
	settingWindow     fyne.Window
	cronJob           = newCron()
	trayIconResource  []byte
)

func main() {
	cronJob.Start()
	setupIcon()
	go startTray()
	logWindow = buildLogWindow()
	settingWindow = buildPrefWindow()
	go Start()
	mainApp.Run()
}

func setupIcon() {
	mainApp.SetIcon(PngIconResource)
	// windows tray icon issue walk around https://github.com/reujab/wallpaper/pull/15
	if runtime.GOOS == "windows" {
		trayIconResource = IcoIconResource.StaticContent
	} else {
		trayIconResource = PngIconResource.StaticContent
	}
}

func buildPrefWindow() fyne.Window {
	settingWindow := mainApp.NewWindow("Preferences")
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
	value := mainApp.Preferences().StringWithFallback("interval", "@daily")
	mainApp.Preferences().SetString("interval", value)
	intervalEntry.SetText(value)
	intervalEntry.SetPlaceHolder("@daily")

	intervalEntry.OnChanged = func(text string) {
		_, err := parseCron(text)
		if err != nil {
			intervalEntryErrorLabel.SetText("Wrong Format")
		} else {
			intervalEntryErrorLabel.SetText("")
			mainApp.Preferences().SetString("interval", text)
			_, err := clearAndSetCron(text)
			checkError("set cron failed", err)
		}
	}

	sortingSelect := getSelect("sorting", sorting)

	firstOrRandomSelect := getSelect("first_or_random", firstOrRandom)

	perferDarkerSelect := getSelect("prefer_darker", preferDarker)

	autorunCheck := widget.NewCheck("autorun", func(toggle bool) {
		mainApp.Preferences().SetBool("autorun", toggle)

		autoStartApp, err := newAutoRun()
		checkError("get autostart app failed", err)
		if toggle {
			autoStartApp.Enable()
		} else {
			autoStartApp.Disable()
		}

	})
	autorunEnabled := mainApp.Preferences().BoolWithFallback("autorun", false)
	autorunCheck.SetChecked(autorunEnabled)

	deepscanCheck := widget.NewCheck("deepscan", func(toggle bool) {
		mainApp.Preferences().SetBool("deepscan", toggle)
	})
	deepscanEnabled := mainApp.Preferences().BoolWithFallback("deepscan", false)
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
func buildLogWindow() fyne.Window {
	logWindow := mainApp.NewWindow("Logs")
	logWindow.SetIcon(PngIconResource)
	logWindow.CenterOnScreen()
	logWindow.Resize(fyne.NewSize(600, 800))
	logWindow.SetCloseIntercept(func() {
		logWindow.Hide()
	})

	logEntry = widget.NewMultiLineEntry()
	logEntry.Disable()

	logWindow.SetContent(container.NewScroll(logEntry))

	logWindow.Hide()
	return logWindow
}
func newLogError(text string, err error) {
	logEntry.Text += time.Now().Format(time.RFC3339) + "\tERROR\t" + text + "\t" + err.Error() + "\n"
}
func newLogInfo(text string) {
	logEntry.Text += time.Now().Format(time.RFC3339) + "\tINFO\t" + text + "\n"
}

func getStringInputBox(name, fallback string) *widget.Entry {
	entry := widget.NewEntry()
	value := mainApp.Preferences().StringWithFallback(name, fallback)
	mainApp.Preferences().SetString(name, value)
	entry.SetText(value)
	entry.SetPlaceHolder(fallback)
	entry.OnChanged = func(text string) {
		mainApp.Preferences().SetString(name, text)
	}
	return entry
}

func getIntInputBox(name string, fallback int, errorMsg *widget.Label) *widget.Entry {
	entry := widget.NewEntry()
	value := mainApp.Preferences().IntWithFallback(name, fallback)
	mainApp.Preferences().SetInt(name, value)
	text := strconv.Itoa(value)
	entry.SetText(text)
	entry.SetPlaceHolder(text)
	entry.OnChanged = func(text string) {
		i, err := strconv.Atoi(text)
		if err != nil {
			errorMsg.SetText("Not a number")
		} else {
			mainApp.Preferences().SetInt(name, i)
			errorMsg.SetText("")
		}
	}
	return entry
}

func getSelect(name string, selection *options) *widget.Select {
	selectEl := widget.NewSelect(selection.getNames(), func(text string) {
		mainApp.Preferences().SetString(name, selection.getValueFromName(text))
	})
	value := mainApp.Preferences().StringWithFallback(name, selection.getValues()[0])
	mainApp.Preferences().SetString(name, value)
	selectEl.SetSelected((*selection)[value])

	return selectEl
}

func checkError(text string, err error) {
	if err != nil {
		newLogError(text, err)
	}
}

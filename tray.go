package main

import "github.com/getlantern/systray"

func onReady() {
	systray.SetIcon(trayIconResource)
	systray.SetTitle("Go Reddit WallPaper")
	systray.SetTooltip("Go Reddit WallPaper")

	mRefresh := systray.AddMenuItem("Refresh Now", "Refresh Now!")
	mPref := systray.AddMenuItem("Preferences", "Change Preferences")
	mLog := systray.AddMenuItem("Logs", "See Logs")
	mQuit := systray.AddMenuItem("Quit", "Quit Go Reddit WallPaper")
	for {
		select {

		case <-mRefresh.ClickedCh:
			go Start()

		case <-mPref.ClickedCh:
			settingWindow.Show()

		case <-mLog.ClickedCh:
			logWindow.Show()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}
func onExit() {
	mainApp.Quit()
}

func startTray() {
	systray.Run(onReady, onExit)
}

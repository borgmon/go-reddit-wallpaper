package main

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/ProtonMail/go-autostart"
)

func getAutorunExec() ([]string, error) {
	dir, err := os.Executable()
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "windows" {
		return []string{dir}, nil
	} else if runtime.GOOS == "linux" {
		return []string{"bash", "-c", dir}, nil
	} else if runtime.GOOS == "darwin" {
		fileName := "~/Library/LaunchAgents/me.borgmon.go-reddit-wallpaper.plist"
		err := ioutil.WriteFile(fileName, plistResource.StaticContent, 0644)
		if err != nil {
			return nil, err
		}
		return []string{"launchctl load " + fileName}, nil
	} else {
		return nil, errors.New("Autorun not implemented")
	}
}

func newAutoRun() (*autostart.App, error) {
	exec, err := getAutorunExec()
	if err != nil {
		return nil, err
	}
	return &autostart.App{
		Name:        "go-reddit-wallpaper",
		DisplayName: "Go Reddit WallPaper",
		Exec:        exec,
	}, nil
}

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/ProtonMail/go-autostart"
)

func getAutorunExec() (error, []string) {
	dir, err := os.Executable()
	if err != nil {
		return err, nil
	}
	if runtime.GOOS == "windows" {
		return nil, []string{dir}
	} else if runtime.GOOS == "linux" {
		return nil, []string{"bash", "-c", dir}
	} else if runtime.GOOS == "darwin" {
		fileName := "~/Library/LaunchAgents/me.borgmon.go-reddit-wallpaper.plist"
		err := ioutil.WriteFile(fileName, PlistResource.StaticContent, 0644)
		if err != nil {
			return err, nil
		}
		return nil, []string{"launchctl load " + fileName}
	} else {
		return errors.New("Autorun not implemented"), nil
	}
}

func newAutoRun() (*autostart.App, error) {
	err, exec := getAutorunExec()
	if err != nil {
		return nil, err
	}
	return &autostart.App{
		Name:        "go-reddit-wallpaper",
		DisplayName: "Go Reddit WallPaper",
		Exec:        exec,
	}, nil
}

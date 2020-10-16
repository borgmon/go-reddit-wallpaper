package main

import "fyne.io/fyne/widget"

func aaa() {
	w1 := MainApp.NewWindow("1")
	w1.CenterOnScreen()
	w1.SetContent(widget.NewVBox(widget.NewButton("hahah", func() {
		w2 := MainApp.NewWindow("1")
		w2.CenterOnScreen()
		w2.SetContent(widget.NewVBox(widget.NewLabel("asd")))
	})))
}

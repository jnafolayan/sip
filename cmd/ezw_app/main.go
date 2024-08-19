package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	win := a.NewWindow("EZW coding")
	win.Resize(fyne.NewSize(800, 640))
	win.SetContent(homeView(win))

	win.ShowAndRun()
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"github.com/jnafolayan/sip/pkg/codec"
)

const WIDTH = 1240
const HEIGHT = 960

func main() {
	a := app.New()
	win := a.NewWindow("EZW coding")
	win.Resize(fyne.NewSize(WIDTH, HEIGHT))
	// win.SetContent(homeView(win))
	win.SetContent(streamView(win, storage.NewFileURI(baseDir+"/lena.png"), codec.DefaultCodecOpts))

	win.ShowAndRun()
}

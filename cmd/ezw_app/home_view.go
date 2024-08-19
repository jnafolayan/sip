package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var baseDir = "/Users/jnafolayan/workspace/projects/sip/resources"

func homeView(win fyne.Window) *fyne.Container {
	fileSelectDialog := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if f == nil {
			return
		}
		fmt.Println(f.URI().String())
	}, win)

	fileSelectDialog.Resize(fyne.NewSize(800, 640))
	fileSelectDialog.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))
	l, _ := storage.ListerForURI(storage.NewFileURI(baseDir))
	fileSelectDialog.SetLocation(l)

	uploadButton := widget.NewButtonWithIcon("Select image to stream", theme.UploadIcon(), func() {
		fileSelectDialog.Show()
	})

	content := container.New(layout.NewCenterLayout(), uploadButton)
	return content
}

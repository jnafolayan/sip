package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jnafolayan/sip/pkg/codec"
)

func streamView(_win fyne.Window, uri fyne.URI, codecOpts codec.CodecOptions) *fyne.Container {
	img := canvas.NewImageFromURI(uri)
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(WIDTH*0.8, 800))
	img.SetMinSize(img.Size())

	stepButton := widget.NewButtonWithIcon("Step", theme.ViewRefreshIcon(), func() {
		// step
	})
	stepButton.Resize(fyne.NewSize(100, 40))

	view := container.New(
		layout.NewVBoxLayout(),
		container.NewHBox(
			layout.NewSpacer(),
			img,
			layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(), stepButton, layout.NewSpacer(),
		),
	)
	view.Resize(fyne.NewSize(WIDTH, HEIGHT))

	return view
}

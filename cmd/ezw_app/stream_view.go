package main

import (
	"bytes"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/ezw"
)

func streamView(win fyne.Window, uri fyne.URI, codecOpts codec.CodecOptions) *fyne.Container {
	img := canvas.NewImageFromURI(uri)
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(WIDTH*0.8, 800))
	img.SetMinSize(img.Size())

	encoder := ezw.NewImageEncoder(codecOpts)
	stepButton := widget.NewButtonWithIcon("Step", theme.ViewRefreshIcon(), func() {
		buf := new(bytes.Buffer)
		encoder.Tick(buf)
		fmt.Println(buf.Len())

		buf.Reset()
	})
	stepButton.Resize(fyne.NewSize(100, 40))
	stepButton.Disable()

	go func() {
		err := encoder.Init(uri.Path())
		if err != nil {
			dialog.ShowError(err, win)
			// FIXME: reload*
			return
		}
		stepButton.Enable()
	}()

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

package main

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/ezw"
)

func streamView(win fyne.Window, uri fyne.URI, codecOpts codec.CodecOptions) *fyne.Container {
	img := canvas.NewImageFromURI(uri)
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(WIDTH*0.8, 800))
	img.SetMinSize(img.Size())

	processing := false

	encoder := ezw.NewImageEncoder(codecOpts)
	decoder := ezw.NewImageDecoder(encoder.SrcSize(), codecOpts)
	stepButton := widget.NewButtonWithIcon("Step", theme.ViewRefreshIcon(), func() {
		if processing {
			return
		}
		processing = true
		buf := new(bytes.Buffer)
		err := encoder.TickJSON(buf)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		err = decoder.DecodeJSONFrame(buf)
		buf.Reset()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		channels := decoder.ReconstructChannels()
		reconstructed := imageutils.ReconstructImageWithAlpha(channels, 255)
		img.Resource = nil
		img.Image = reconstructed
		img.File = ""
		img.Refresh()
		processing = false

		// debug.DrawSignal2D(channels[0], image.Rect(0, 0, 15, 15), "temp.jpg")

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

		err = decoder.Init()
		if err != nil {
			dialog.ShowError(err, win)
			// FIXME: reload*
			return
		}

		decoder.SetDestSize(encoder.SrcSize())

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
